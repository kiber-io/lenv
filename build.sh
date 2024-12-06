#!/bin/bash

OUTPUT_DIR="./build"

mkdir -p "$OUTPUT_DIR"

build() {
    local os=$1
    local arch=$2
    local output_name=$3

    echo "Building for $os $arch..."

    GOOS=$os GOARCH=$arch go build -o "$OUTPUT_DIR/$output_name" ./cmd/lenv

    if [ $? -eq 0 ]; then
        echo "Successfully built $output_name"
    else
        echo "Failed to build $output_name"
    fi
}

build "windows" "amd64" "lenv_win_x64.exe"
build "linux" "amd64" "lenv_linux_x64"
build "linux" "arm64" "lenv_linux_arm64"

echo "Build process completed."
