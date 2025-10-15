#!/bin/bash
set -e

echo "Quick Deploy Script for ts-unduck"
echo "=================================="
echo ""

INSTALL_DIR="$HOME/unduck"

if [ ! -d "$INSTALL_DIR" ]; then
    echo "Cloning ts-unduck..."
    git clone https://github.com/j4ckxyz/ts-unduck.git "$INSTALL_DIR"
    cd "$INSTALL_DIR"
else
    echo "Updating ts-unduck..."
    cd "$INSTALL_DIR"
    git pull
fi

echo "Running installer..."
./install.sh
