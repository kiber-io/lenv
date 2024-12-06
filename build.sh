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

build "windows" "amd64" "lenv_win64.exe"
build "windows" "386" "lenv_win32.exe"
build "linux" "amd64" "lenv_linux64"
build "linux" "386" "lenv_linux32"

echo "Build process completed."
