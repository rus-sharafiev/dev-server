#!/usr/bin/env pwsh

$ErrorActionPreference = 'Stop'

$BinDir = "${Home}\.dev"

$Zip = "$BinDir\dev.zip"
$Exe = "$BinDir\dev.exe"

$DownloadUrl = "https://github.com/rus-sharafiev/dev-server/releases/latest/download/dev.zip"

if (!(Test-Path $BinDir)) {
    New-Item $BinDir -ItemType Directory | Out-Null
}

curl.exe -Lo $Zip $DownloadUrl

tar.exe xf $Exe -C $BinDir

Remove-Item $Zip

$User = [System.EnvironmentVariableTarget]::User
$Path = [System.Environment]::GetEnvironmentVariable('Path', $User)
if (!(";${Path};".ToLower() -like "*;${BinDir};*".ToLower())) {
    [System.Environment]::SetEnvironmentVariable('Path', "${Path};${BinDir}", $User)
    $Env:Path += ";${BinDir}"
}

Write-Output "Dev-server has been installed successfully to ${Exe}"