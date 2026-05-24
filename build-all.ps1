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
  
  if ($goos -eq "windows") {
    $ext = ".exe"
  } elseif ($goos -eq "js" -or $goos -eq "wasip1") {
    $ext = ".wasm"
  } else {
    $ext = ""
  }
  
  $binaryName = "easyedit-$version-$goos-$goarch$ext"
  $outputPath = "$outputDir\$binaryName"
  
  Write-Host "Building $goos/$goarch ... " -NoNewline
  $env:GOOS = $goos
  $env:GOARCH = $goarch
  $env:CGO_ENABLED = 0
  
  $result = & go build -ldflags="-s -w" -o "$outputPath" 2>&1
  if ($LASTEXITCODE -eq 0) {
    $size = (Get-Item "$outputPath").Length
    $sizeKB = [math]::Round($size / 1KB, 1)
    Write-Host "OK ($sizeKB KB)" -ForegroundColor Green
    $built += @{Path = $outputPath; Name = $binaryName}
  } else {
    Write-Host "FAILED" -ForegroundColor Red
    Write-Host "  $result"
    $failed += $p
  }
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
