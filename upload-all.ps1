$release = "v0.2.0"
$buildDir = "D:\easyedit-builds"
$repo = "0xA672/Easyedit"

Write-Host "Uploading all assets to $release ..."

Get-ChildItem $buildDir -File | Where-Object { $_.Extension -ne ".txt" } | ForEach-Object {
  $path = $_.FullName
  $name = $_.Name
  Write-Host "Uploading $name ... " -NoNewline
  
  $result = gh release upload $release "`"$path`"#`"$name`"" --repo $repo --clobber 2>&1
  if ($LASTEXITCODE -eq 0) {
    Write-Host "OK" -ForegroundColor Green
  } else {
    Write-Host "FAILED: $result" -ForegroundColor Red
  }
}

Write-Host "`nDone!"
