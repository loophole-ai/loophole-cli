$ErrorActionPreference = 'Stop'

$AppName = "loophole"
$RepoOwner = "loophole-ai"
$RepoName = "loophole-cli"

function Get-LatestVersion {
    $url = "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"
    $response = Invoke-RestMethod -Uri $url
    return $response.tag_name.TrimStart('v')
}

$Version = if ($env:VERSION) { $env:VERSION } else { Get-LatestVersion }

$OS = "windows"
$Arch = if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") { "amd64" } else { "arm64" }

$FileName = "$AppName-$OS-$Arch.exe"
$DownloadUrl = "https://github.com/$RepoOwner/$RepoName/releases/download/v$Version/$FileName"

$InstallDir = Join-Path $HOME ".loophole\bin"
if (!(Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

$DestFile = Join-Path $InstallDir "$AppName.exe"

Write-Host "Downloading $AppName v$Version..." -ForegroundColor Cyan
Write-Host "URL: $DownloadUrl" -ForegroundColor Gray

try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $DestFile
} catch {
    Write-Error "Failed to download $AppName. Please check if version v$Version exists."
    exit 1
}

Write-Host "Successfully installed $AppName to $InstallDir" -ForegroundColor Green

# Add to PATH if not already there
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "Adding $InstallDir to User PATH..." -ForegroundColor Yellow
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path += ";$InstallDir"
    Write-Host "Please restart your terminal or run: `$env:Path += ';$InstallDir'`" -ForegroundColor Cyan
} else {
    Write-Host "$AppName is already in your PATH." -ForegroundColor Gray
}

Write-Host "Installation complete! Try running: loophole --version" -ForegroundColor Green
