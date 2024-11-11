#!/bin/bash

set -e

ENV_JAVAENV_HOME="JAVAENV_HOME"
ENV_JAVA_HOME="JAVA_HOME"
ENV_PATH="PATH"

javaenvHomePath="$HOME/.javaenv"
javaenvHomeBinPath="$javaenvHomePath/bin"
currentJdkPath="$javaenvHomePath/currentjdk"

if [ "$(id -u)" -eq 0 ]; then
    echo "Error: Please run this script as a non-administrator."
    exit 1
fi

if [ "$(uname -m)" = "x86_64" ]; then
    assetName="javaenv_linux64"
else
    assetName="javaenv_linux32"
fi

latestRelease=$(curl -s "https://api.github.com/repos/kiber-io/javaenv/releases/latest")
latestAsset=$(echo "$latestRelease" | grep -oP "(?<=\"browser_download_url\": \").*${assetName}\"" | tr -d '\"')

if [ -z "$latestAsset" ]; then
    echo "Error: Could not find the asset for $assetName in the latest release."
    exit 1
fi

mkdir -p "$javaenvHomeBinPath"
mkdir -p "$currentJdkPath"

downloadPath="$javaenvHomeBinPath/javaenv"
curl -L "$latestAsset" -o "$downloadPath"
chmod +x "$downloadPath"

path=$(echo "$PATH")
if [[ ! "$path" == *"$ENV_JAVA_HOME/bin"* ]]; then
    path="$javaenvHomePath/currentjdk/bin:$path"
fi
if [[ ! "$path" == *"$ENV_JAVAENV_HOME/bin"* ]]; then
    path="$javaenvHomeBinPath:$path"
fi

export JAVAENV_HOME="$javaenvHomePath"
export JAVA_HOME="$currentJdkPath"
export PATH="$path"

echo "export JAVAENV_HOME=\"$javaenvHomePath\"" >> ~/.bashrc
echo "export JAVA_HOME=\"$currentJdkPath\"" >> ~/.bashrc
echo "export PATH=\"$javaenvHomeBinPath:\$JAVA_HOME/bin:\$PATH\"" >> ~/.bashrc
echo "export JAVA_HOME=\"$currentJdkPath\"" >> ~/.profile

echo "Installation completed. Please restart your terminal to start using javaenv."
