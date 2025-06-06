#!/bin/bash

# build project first
uv build

# find wheel path in dist folder, set the name and path
export PYAPP_PROJECT_NAME="dosh"
export PYAPP_PROJECT_PATH="$(find $(pwd)/dist -name "dosh*.whl" -type f | head -1)"

# download pyapp source code, build it, and rename the binary
curl https://github.com/ofek/pyapp/releases/latest/download/source.tar.gz -Lo pyapp-source.tar.gz
tar -xzf pyapp-source.tar.gz
mv pyapp-v* pyapp-latest
cd pyapp-latest
cargo build --release
mkdir -p ../bin
mv target/release/pyapp ../bin/dosh-cli-linux && chmod +x ../bin/dosh-cli-linux
