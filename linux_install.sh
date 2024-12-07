#!/bin/bash

DEBUG=false

while getopts "d" opt; do
  case $opt in
    d) DEBUG=true ;;
    *) echo "Invalid option"; exit 1 ;;
  esac
done

initialize_environment_variables() {
  lenv_home_path="$HOME/.lenv"
}

test_admin() {
  if [ "$(id -u)" -eq 0 ]; then
    echo "Error: Please run this script as a non-administrator."
    exit 1
  fi
}

get_asset_name() {
  os=$(uname -o)
  arch=$(uname -m)

  case "$os" in
    GNU/Linux) os="linux" ;;
    Android) os="android" ;;
    *)
      echo "Error: Unknown operating system $os"
      exit 1
      ;;
  esac

  case "$arch" in
    x86_64) arch="x64" ;;
    aarch64) arch="arm64" ;;
    *)
      echo "Error: Unknown architecture $arch"
      exit 1
      ;;
  esac

  echo "lenv_${os}_${arch}"
}

get_latest_asset() {
  latest_release=$(curl -s https://api.github.com/repos/kiber-io/lenv/releases/latest)
  asset_name=$(get_asset_name)
  download_url=$(echo "$latest_release" | grep -oP "(?<=\"browser_download_url\": \").*${asset_name}.*(?=\")")

  if [ -z "$download_url" ]; then
    echo "Error: Could not find the asset for $asset_name in the latest release."
    exit 1
  fi
  echo "$download_url"
}

new_directories() {
  mkdir -p "$lenv_home_path/bin"
  if [ ! -e "$lenv_home_path/java/current" ] && [ ! -L "$lenv_home_path/java/current" ]; then
    mkdir -p "$lenv_home_path/java/current"
  fi
  if [ ! -d "$lenv_home_path/python/current" ] && [ ! -L "$lenv_home_path/python/current" ]; then
    mkdir -p "$lenv_home_path/python/current"
  fi
}

get_asset() {
  if [ "$DEBUG" = true ]; then
    local_file_path="$PWD/build/$(get_asset_name)"
    cp "$local_file_path" "$lenv_home_path/bin/lenv"
  else
    download_url=$(get_latest_asset)
    curl -L "$download_url" -o "$lenv_home_path/bin/lenv"
  fi
  chmod +x "$lenv_home_path/bin/lenv"
}

update_environment_variables() {
  profile_file="$HOME/.profile"
  bashrc_file="$HOME/.bashrc"

  ensure_newline() {
    file=$1
    [ -s "$file" ] && tail -c1 "$file" | read -r _ || echo >> "$file"
  }

  ensure_newline "$profile_file"
  if ! grep -q "export LENV_HOME=$lenv_home_path" "$profile_file"; then
    echo "export LENV_HOME=$lenv_home_path" >> "$profile_file"
  fi
  if ! grep -q "export JAVA_HOME=\$LENV_HOME/java/current" "$profile_file"; then
    echo "export JAVA_HOME=\$LENV_HOME/java/current" >> "$profile_file"
  fi
  if ! grep -q "export PATH=\$LENV_HOME/bin:\$JAVA_HOME/bin:\$PATH" "$profile_file"; then
    echo "export PATH=\$LENV_HOME/bin:\$JAVA_HOME/bin:\$PATH" >> "$profile_file"
  fi

  ensure_newline "$bashrc_file"
  if ! grep -q "export LENV_HOME=$lenv_home_path" "$bashrc_file"; then
    echo "export LENV_HOME=$lenv_home_path" >> "$bashrc_file"
  fi
  if ! grep -q "export PATH=\$LENV_HOME/bin:\$JAVA_HOME/bin:\$PATH" "$bashrc_file"; then
    echo "export PATH=\$LENV_HOME/bin:\$JAVA_HOME/bin:\$PATH" >> "$bashrc_file"
  fi
  if ! grep -q "export PATH=\$LENV_HOME/python/current/bin:\$PATH" "$bashrc_file"; then
    echo "export PATH=\$LENV_HOME/python/current/bin:\$PATH" >> "$bashrc_file"
  fi
}

main() {
  test_admin
  initialize_environment_variables
  new_directories
  get_asset
  update_environment_variables
  echo "Installation completed. Please restart your terminal to start using lenv."
}

main
