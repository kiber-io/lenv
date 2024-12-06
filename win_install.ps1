param (
    [switch]$Debug
)

function Initialize-EnvironmentVariables {
    $lenvHomePath = "$env:USERPROFILE\.lenv"
    $envVars = @{
        ENV_LENV_HOME   = "LENV_HOME"
        ENV_JAVA_HOME   = "JAVA_HOME"
        ENV_PATH        = "Path"
        lenvHomePath    = $lenvHomePath
        lenvHomeBinPath = "$lenvHomePath\bin"
        javaCurrentPath = "$lenvHomePath\java\current"
    }
    return $envVars
}

function Test-Admin {
    $isAdmin = ([Security.Principal.WindowsIdentity]::GetCurrent()).Groups -match 'S-1-5-32-544'
    if ($isAdmin) {
        Write-Output "Error: Please run this script as a non-administrator."
        exit
    }
}

function Get-AssetName {
    if ([System.Environment]::Is64BitOperatingSystem) {
        return "lenv_win_x64.exe"
    }
    else {
        Write-Output "Error: unsupported architecture."
        exit 1
    }
}

function Get-LatestAsset {
    $latestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/kiber-io/lenv/releases/latest"
    $assetName = Get-AssetName
    $latestAsset = $latestRelease.assets | Where-Object { $_.name -eq $assetName }

    if (-not $latestAsset) {
        Write-Output "Error: Could not find the asset for $assetName in the latest release."
        exit 1
    }
    return $latestAsset.browser_download_url
}

function New-Directories ($envVars) {
    if (!(Test-Path -Path $envVars.lenvHomeBinPath)) {
        New-Item -ItemType Directory -Path $envVars.lenvHomeBinPath | Out-Null
    }
    if (!(Test-Path -Path $envVars.javaCurrentPath)) {
        New-Item -ItemType Directory -Path $envVars.javaCurrentPath | Out-Null
    }
}

function Get-Asset ($envVars) {
    if ($Debug) {
        $localFilePath = ""
        $localFilePath = "$PSScriptRoot\build\" + (Get-AssetName)
        $destinationPath = "$($envVars.lenvHomeBinPath)\lenv.exe"
        Copy-Item -Path $localFilePath -Destination $destinationPath
    } else {
        $downloadUrl = Get-LatestAsset
        $downloadPath = "$($envVars.lenvHomeBinPath)\lenv.exe"
        Invoke-WebRequest -Uri $downloadUrl -OutFile $downloadPath
    }
}

function Update-EnvironmentVariables ($envVars) {
    $path = [System.Environment]::GetEnvironmentVariable($envVars.ENV_PATH, [System.EnvironmentVariableTarget]::User)
    $javaHomePath = [System.Environment]::GetEnvironmentVariable($envVars.ENV_JAVA_HOME, [System.EnvironmentVariableTarget]::User)
    if ($null -eq $javaHomePath -or -not $path.Contains("%$($envVars.ENV_JAVA_HOME)%\bin")) {
        $path = "%$($envVars.ENV_JAVA_HOME)%\bin;$path"
    }

    $lenvPath = [System.Environment]::GetEnvironmentVariable($envVars.ENV_LENV_HOME, [System.EnvironmentVariableTarget]::User)
    if ($null -eq $lenvPath -or -not $path.Contains("%$($envVars.ENV_LENV_HOME)%\bin")) {
        $path = "%$($envVars.ENV_LENV_HOME)%\bin;$path"
    }

    [System.Environment]::SetEnvironmentVariable($envVars.ENV_LENV_HOME, $envVars.lenvHomePath, [System.EnvironmentVariableTarget]::User)
    [System.Environment]::SetEnvironmentVariable($envVars.ENV_JAVA_HOME, "%$($envVars.ENV_LENV_HOME)%\java\current", [System.EnvironmentVariableTarget]::User)
    [System.Environment]::SetEnvironmentVariable($envVars.ENV_PATH, $path, [System.EnvironmentVariableTarget]::User)
}

function Main {
    Test-Admin
    $envVars = Initialize-EnvironmentVariables
    New-Directories $envVars
    Get-Asset $envVars
    Update-EnvironmentVariables $envVars
    Write-Output "Installation completed. Please restart your terminal to start using lenv."
}

Main
