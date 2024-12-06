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

Build "windows" "amd64" "lenv_win_x64.exe"
Build "linux" "amd64" "lenv_linux_x64"

Write-Output "Build process completed."
