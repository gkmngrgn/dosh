#!/bin/bash

# Set output binary name (can be overridden by command line argument or environment variable)
OUTPUT_BINARY="${DOSH_BINARY_NAME:-dosh}"

# find wheel path in dist folder, set the name and path
export PYAPP_PROJECT_NAME="dosh"
export PYAPP_PROJECT_PATH="$(find $(pwd)/dist -name "dosh*.whl" -type f | head -1)"

# download pyapp source code, build it
curl https://github.com/ofek/pyapp/releases/latest/download/source.tar.gz -Lo pyapp-source.tar.gz
tar -xzf pyapp-source.tar.gz
mv pyapp-v* pyapp-latest
cd pyapp-latest
cargo build --release

# and rename the binary
mkdir -p ../bin
mv target/release/pyapp "../bin/${OUTPUT_BINARY}"
chmod +x "../bin/${OUTPUT_BINARY}"
