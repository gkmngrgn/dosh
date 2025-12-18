#!/bin/sh
set -e

# Default installation directory
install_dir="$HOME/.local/bin"

# Parse command-line arguments
while [ $# -gt 0 ]; do
  case "$1" in
    --install-dir)
      if [ -z "$2" ]; then
        echo "Error: --install-dir requires a value" >&2
        exit 1
      fi
      install_dir="$2"
      shift # past argument
      shift # past value
      ;;
    *)
      # unknown option
      shift # past argument
      ;;
  esac
done

os=$(uname -s | tr '[:upper:]' '[:lower:]')
# Map darwin to macos for release asset names
if [ "$os" = "darwin" ]; then
    os="macos"
fi
architecture=$(uname -m)
download_url="https://github.com/gkmngrgn/dosh/releases/latest/download/dosh-$os-$architecture"
temp_dir=$(mktemp -d)
bin_file="$install_dir/dosh"

echo "Operating System: $os"
echo "Architecture: $architecture"
echo "Temporary directory: $temp_dir"
echo "Download URL: $download_url"
echo "Installation directory: $install_dir"

printf "\nSTEP 1: Downloading DOSH...\n"
curl -L "$download_url" -o "$temp_dir/dosh"

printf "\nSTEP 2: Installing DOSH CLI to %s...\n" "$bin_file"
if [ -f "$bin_file" ]; then
    mv "$bin_file" "$temp_dir/dosh.old"
else
    # make sure if local bin folder exists
    mkdir -p "$install_dir"
fi

mv "$temp_dir/dosh" "$bin_file"
chmod +x "$bin_file"

# Expand install_dir to absolute path for PATH comparison
# This handles cases where user provides ~/bin but PATH has /home/user/bin
expanded_install_dir=$(cd "$install_dir" 2>/dev/null && pwd) || expanded_install_dir="$install_dir"

if ! echo ":$PATH:" | grep -q ":$expanded_install_dir:" && ! echo ":$PATH:" | grep -q ":$install_dir:"; then
    printf "\n\033[1;33mWARNING: '%s' is not in your PATH.\033[0m" "$install_dir"
    printf "\n\033[1;33mYou should add the following line to your shell configuration file (e.g., ~/.bashrc, ~/.zshrc):\033[0m\n"
    printf '\n\033[1;33m  export PATH="%s:$PATH"\033[0m\n' "$install_dir"
fi

printf "\nSTEP 3: Done! You can delete the temporary directory if you want:"
printf "\n%s\n" "$temp_dir"
