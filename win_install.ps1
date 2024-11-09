$repo = "kiber-io/javaenv"
$rootDir = "$env:USERPROFILE\.javaenv"
$installDir = "$rootDir\bin"
$envVarPath = "$installDir\javaenv.exe"
$currentJdkPath = "$rootDir\currentjdk\bin"

$assetName = if ([System.Environment]::Is64BitOperatingSystem) {
    "javaenv_win64.exe"
} else {
    "javaenv_win32.exe"
}

$latestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest"
$latestAsset = $latestRelease.assets | Where-Object { $_.name -eq $assetName }

if (-not $latestAsset) {
    Write-Output "Error: Could not find the asset for $assetName in the latest release."
    exit 1
}

if (!(Test-Path -Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

$downloadPath = "$installDir\javaenv.exe"
Invoke-WebRequest -Uri $latestAsset.browser_download_url -OutFile $downloadPath

$currentPath = [System.Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::User)
if (-not $currentPath.Contains($installDir)) {
    $currentPath = "$currentPath;$installDir"
}
if (-not $currentPath.Contains($currentJdkPath)) {
    $currentPath = "$currentPath;$currentJdkPath"
}
[System.Environment]::SetEnvironmentVariable("Path", $currentPath, [System.EnvironmentVariableTarget]::User)

Write-Output "Installation completed. You can now use 'javaenv' from the command line. Please restart your command line interface."
