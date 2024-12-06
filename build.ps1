$OUTPUT_DIR = "./build"

if (-Not (Test-Path -Path $OUTPUT_DIR)) {
    New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null
}

function Build {
    param (
        [string]$os,
        [string]$arch,
        [string]$output_name
    )

    Write-Output "Building for $os $arch..."

    $env:GOOS = $os
    $env:GOARCH = $arch
    & go build -o "$OUTPUT_DIR\$output_name" ./cmd/lenv

    if ($?) {
        Write-Output "Successfully built $output_name"
    } else {
        Write-Output "Failed to build $output_name"
    }
}

Build "windows" "amd64" "lenv_win64.exe"
Build "windows" "386" "lenv_win32.exe"
Build "linux" "amd64" "lenv_linux64"
Build "linux" "386" "lenv_linux32"

Write-Output "Build process completed."
