$ENV_JAVAENV_HOME = "JAVAENV_HOME"
$ENV_JAVA_HOME = "JAVA_HOME"
$ENV_PATH = "Path"

$javaenvHomePath = "$env:USERPROFILE\.javaenv"
$javaenvHomeBinPath = "$javaenvHomePath\bin"
$currentJdkPath = "$javaenvHomePath\currentjdk"

$isAdmin = ([Security.Principal.WindowsIdentity]::GetCurrent()).Groups -match 'S-1-5-32-544'

if ($isAdmin) {
    Write-Output "Error: Please run this script as a non-administrator."
    exit
}

$assetName = if ([System.Environment]::Is64BitOperatingSystem) {
    "javaenv_win64.exe"
} else {
    "javaenv_win32.exe"
}

$latestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/kiber-io/javaenv/releases/latest"
$latestAsset = $latestRelease.assets | Where-Object { $_.name -eq $assetName }

if (-not $latestAsset) {
    Write-Output "Error: Could not find the asset for $assetName in the latest release."
    exit 1
}

# create directories only when asset is found to avoid creating unnecessary directories
if (!(Test-Path -Path $javaenvHomeBinPath)) {
    New-Item -ItemType Directory -Path $javaenvHomeBinPath | Out-Null
}
if (!(Test-Path -Path $currentJdkPath)) {
    New-Item -ItemType Directory -Path $currentJdkPath | Out-Null
}

$downloadPath = "$javaenvHomeBinPath\javaenv.exe"
Invoke-WebRequest -Uri $latestAsset.browser_download_url -OutFile $downloadPath

$path = [System.Environment]::GetEnvironmentVariable($ENV_PATH, [System.EnvironmentVariableTarget]::User)
$javaHomePath = [System.Environment]::GetEnvironmentVariable($ENV_JAVA_HOME, [System.EnvironmentVariableTarget]::User)
# if JAVA_HOME\bin is not in the path, add it
if ($null -eq $javaHomePath -or -not $path.Contains("%$ENV_JAVA_HOME%\bin")) {
    $path = "%$ENV_JAVA_HOME%\bin;$path"
}

$javaenvPath = [System.Environment]::GetEnvironmentVariable($ENV_JAVAENV_HOME, [System.EnvironmentVariableTarget]::User)
# if JAVAENV_HOME\bin is not set or it is not in the path, add it
if ($null -eq $javaenvPath -or -not $path.Contains("%$ENV_JAVAENV_HOME%\bin")) {
    $path = "%$ENV_JAVAENV_HOME%\bin;$path"
}
[System.Environment]::SetEnvironmentVariable($ENV_JAVAENV_HOME, $javaenvHomePath, [System.EnvironmentVariableTarget]::User)
[System.Environment]::SetEnvironmentVariable($ENV_JAVA_HOME, "%$ENV_JAVAENV_HOME%\currentjdk", [System.EnvironmentVariableTarget]::User)
[System.Environment]::SetEnvironmentVariable($ENV_PATH, $path, [System.EnvironmentVariableTarget]::User)

Write-Output "Installation completed. Please restart your terminal to start using javaenv."
