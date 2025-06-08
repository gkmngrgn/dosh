#!/bin/sh
set -e

os=$(uname -s | tr '[:upper:]' '[:lower:]')
architecture=$(uname -m)
download_url="https://github.com/gkmngrgn/dosh/releases/latest/download/dosh-$os-$architecture"
temp_dir=$(mktemp -d)
local_dir="$HOME/.local"
bin_file="$local_dir/bin/dosh"

echo "Operating System: $os"
echo "Architecture: $architecture"
echo "Temporary directory: $temp_dir"
echo "Download URL: $download_url"

printf "\nSTEP 1: Downloading DOSH...\n"
curl -L "$download_url" -o "$temp_dir/dosh"

printf "\nSTEP 2: Installing DOSH CLI...\n"
if [ -f "$bin_file" ]; then
    mv "$bin_file" "$temp_dir/dosh.old"
else
    # make sure if local bin folder exists
    mkdir -p "$local_dir/bin"
fi

mv "$temp_dir/dosh" "$bin_file"
chmod +x "$bin_file"

printf "\nSTEP 3: Done! You can delete the temporary directory if you want:"
printf "\n%s\n" "$temp_dir"
