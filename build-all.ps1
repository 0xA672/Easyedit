$version = "v0.2.0"
$outputDir = "D:\easyedit-builds"
New-Item -ItemType Directory -Force -Path $outputDir | Out-Null

# Use D: drive for Go build cache to avoid C: disk full
$env:GOCACHE = "D:\go-cache"
$env:GOTMPDIR = "D:\go-tmp"
New-Item -ItemType Directory -Force -Path $env:GOCACHE | Out-Null
New-Item -ItemType Directory -Force -Path $env:GOTMPDIR | Out-Null

$platforms = @(
  "aix/ppc64"
  "android/386"
  "android/amd64"
  "android/arm"
  "android/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "dragonfly/amd64"
  "freebsd/386"
  "freebsd/amd64"
  "freebsd/arm"
  "freebsd/arm64"
  "freebsd/riscv64"
  "illumos/amd64"
  "ios/amd64"
  "ios/arm64"
  "js/wasm"
  "linux/386"
  "linux/amd64"
  "linux/arm"
  "linux/arm64"
  "linux/loong64"
  "linux/mips"
  "linux/mips64"
  "linux/mips64le"
  "linux/mipsle"
  "linux/ppc64"
  "linux/ppc64le"
  "linux/riscv64"
  "linux/s390x"
  "netbsd/386"
  "netbsd/amd64"
  "netbsd/arm"
  "netbsd/arm64"
  "openbsd/386"
  "openbsd/amd64"
  "openbsd/arm"
  "openbsd/arm64"
  "openbsd/ppc64"
  "openbsd/riscv64"
  "plan9/386"
  "plan9/amd64"
  "plan9/arm"
  "solaris/amd64"
  "wasip1/wasm"
  "windows/386"
  "windows/amd64"
  "windows/arm64"
)

$built = @()
$failed = @()

foreach ($p in $platforms) {
  $parts = $p.Split("/")
  $goos = $parts[0]
  $goarch = $parts[1]

  # Determine binary extension (no version/platform in binary name)
  if ($goos -eq "windows") {
    $binaryName = "easyedit.exe"
  } elseif ($goos -eq "js" -or $goos -eq "wasip1") {
    $binaryName = "easyedit.wasm"
  } else {
    $binaryName = "easyedit"
  }

  # Archive name: easyedit-v{version}-{os}-{arch}.zip|tar.gz
  if ($goos -eq "windows") {
    $archiveName = "easyedit-$version-$goos-$goarch.zip"
  } elseif ($goos -eq "js" -or $goos -eq "wasip1") {
    # .wasm files are single artifacts, no archive needed
    $archiveName = "easyedit-$version-$goos-$goarch.wasm"
  } else {
    $archiveName = "easyedit-$version-$goos-$goarch.tar.gz"
  }

  $outputPath = "$outputDir\$archiveName"
  $tmpDir = "$outputDir\tmp-$goos-$goarch"

  Write-Host "Building $goos/$goarch ... " -NoNewline
  $env:GOOS = $goos
  $env:GOARCH = $goarch
  $env:CGO_ENABLED = 0

  # Create temp dir and build binary there with simple name
  New-Item -ItemType Directory -Force -Path $tmpDir | Out-Null
  $binPath = "$tmpDir\$binaryName"

  $result = & go build -ldflags="-s -w" -o "$binPath" 2>&1
  if ($LASTEXITCODE -eq 0) {
    # Package into archive
    if ($goos -eq "js" -or $goos -eq "wasip1") {
      # Just rename the file to the archive name (single .wasm artifact)
      Move-Item -Path "$binPath" -Destination "$outputPath" -Force
    } elseif ($goos -eq "windows") {
      # ZIP
      Compress-Archive -Path "$tmpDir\*" -DestinationPath "$outputPath" -Force
    } else {
      # tar.gz
      tar -czf "$outputPath" -C "$tmpDir" "$binaryName" 2>$null
    }

    if (Test-Path "$outputPath") {
      $size = (Get-Item "$outputPath").Length
      $sizeKB = [math]::Round($size / 1KB, 1)
      Write-Host "OK ($sizeKB KB)" -ForegroundColor Green
      $built += @{Path = $outputPath; Name = $archiveName}
    } else {
      Write-Host "FAILED (packaging)" -ForegroundColor Red
      $failed += $p
    }
  } else {
    Write-Host "FAILED" -ForegroundColor Red
    Write-Host "  $result"
    $failed += $p
  }

  # Cleanup temp dir
  Remove-Item -Path $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
}

Write-Host "`n=== Summary ==="
Write-Host "Built: $($built.Count)" -ForegroundColor Green
Write-Host "Failed: $($failed.Count)" -ForegroundColor Red
if ($failed.Count -gt 0) {
  Write-Host "Failed platforms: $($failed -join ', ')" -ForegroundColor Red
}

# Save list for upload
$built | ForEach-Object { $_.Path } | Set-Content "$outputDir\.built-list.txt"
$built | ForEach-Object { "$($_.Name):$($_.Path)" } | Set-Content "$outputDir\.assets.txt"

Write-Host "`nBuild output: $outputDir"
