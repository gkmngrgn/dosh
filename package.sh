#!/bin/bash
set -e

export PATH="$HOME/.local/bin:$HOME/bin:$PATH"

if which uv &> /dev/null
then
    echo "uv found in path."
else
    curl -LsSf https://astral.sh/uv/install.sh | sh
fi

# install python and project dependencies first
uv sync --no-dev --group package

OS_NAME=$(uv run python -c 'import platform; print(platform.system().lower())')
ARCH_TYPE=$(uv run python -c 'import platform; print(platform.machine().lower())')
PY_VERSION=$(uv run python -c 'import sys; print(".".join(map(str, sys.version_info[:3])))')
DOSH_VERSION=$(uv run python -m dosh.cli version)
DIR_NAME="dosh-cli-${OS_NAME}-${ARCH_TYPE}-${DOSH_VERSION}"

echo "---"
echo "PYTHON PATH    : $(uv run which python)"
echo "PYTHON VERSION : ${PY_VERSION}"
echo "DIRECTORY      : ${DIR_NAME}"
echo "---"

uv run pyinstaller dosh/cli.py \
    --name=dosh \
    --copy-metadata=dosh \
    --console \
    --noconfirm \
    --clean \
    --additional-hooks-dir=pyinstaller_hooks
[ -d "${DIR_NAME}" ] && rm -rf "${DIR_NAME}"

mv dist/dosh "${DIR_NAME}"
