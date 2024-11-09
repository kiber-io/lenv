#!/bin/bash

repo="kiber-io/javaenv"
rootDir="$HOME/.javaenv"
binDir="$rootDir/bin"
envVarPath="$binDir/javaenv"
currentJdkPath="$rootDir/currentjdk/bin"

if [[ $(uname -m) == "x86_64" ]]; then
    assetName="javaenv_linux64"
else
    assetName="javaenv_linux32"
fi

latestRelease=$(curl -s "https://api.github.com/repos/$repo/releases/latest")
latestAsset=$(echo "$latestRelease" | grep -oP "(?<=\"browser_download_url\": \").*${assetName}\"" | tr -d '\"')

if [[ -z $latestAsset ]]; then
    echo "Error: Could not find the asset for $assetName in the latest release."
    exit 1
fi

mkdir -p "$binDir"

downloadPath="$envVarPath"
curl -L "$latestAsset" -o "$downloadPath"
chmod +x "$downloadPath"

if ! echo "$PATH" | grep -q "$binDir"; then
    echo "export PATH=\"$binDir:\$PATH\"" >> ~/.bashrc
fi

if ! echo "$PATH" | grep -q "$currentJdkPath"; then
    echo "export PATH=\"$currentJdkPath:\$PATH\"" >> ~/.bashrc
fi

echo "Installation completed. Please restart your terminal or run 'source ~/.bashrc' to use 'javaenv' from the command line."
