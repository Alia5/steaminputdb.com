#!/usr/bin/env sh

set -e

REPO="Alia5/steaminputdb.com"
API_URL="https://api.github.com/repos/${REPO}/releases/latest"

echo "Fetching latest SteamInputDB Buddy release..."
RELEASE_DATA=$(curl -fsSL "$API_URL")
VERSION=$(printf '%s' "$RELEASE_DATA" \
    | grep -Eo '"tag_name"[[:space:]]*:[[:space:]]*"[^"]+"' \
    | head -n 1 \
| cut -d'"' -f4)

if [ -z "$VERSION" ]; then
    echo "Error: Could not fetch release info"
    exit 1
fi

echo "Version: $VERSION"

ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    *)
        echo "Error: Unsupported architecture: $ARCH"
        echo "Only x86_64 (amd64) is supported."
        exit 1
    ;;
esac

BINARY_NAME="steaminputdb-buddy-linux-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"

INSTALL_DIR="$HOME/.local/bin"
INSTALL_PATH="$INSTALL_DIR/steaminputdb-buddy"

get_version() {
    "$1" --help 2>/dev/null | grep -Eo 'SteamInputDB Buddy - v[^ ]+' | head -1 | sed 's/SteamInputDB Buddy - v//'
}

IS_UPDATE=0
SKIP_INSTALL=0
OLD_VERSION=""
if [ -f "$INSTALL_PATH" ]; then
    IS_UPDATE=1
    echo "Existing installation detected at $INSTALL_PATH"
    
    OLD_VERSION=$(get_version "$INSTALL_PATH")
    if [ -z "$OLD_VERSION" ]; then OLD_VERSION="unknown"; fi
    echo "Installed version: $OLD_VERSION"
fi

echo "Downloading from: $DOWNLOAD_URL"
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

TEMP_BIN="$TEMP_DIR/steaminputdb-buddy"
if ! curl -fsSL -o "$TEMP_BIN" "$DOWNLOAD_URL"; then
    echo "Error: Could not download SteamInputDB Buddy binary"
    exit 1
fi

chmod +x "$TEMP_BIN"

NEW_VERSION=$(get_version "$TEMP_BIN")
if [ -z "$NEW_VERSION" ]; then NEW_VERSION="unknown"; fi
echo "Downloaded version: $NEW_VERSION"

if [ "$IS_UPDATE" -eq 1 ]; then
    if [ "$NEW_VERSION" = "$OLD_VERSION" ] && [ "$NEW_VERSION" != "unknown" ]; then
        echo "Already at latest version. Skipping."
        SKIP_INSTALL=1
    fi
fi

if [ "$SKIP_INSTALL" -eq 0 ]; then
    mkdir -p "$INSTALL_DIR"
    
    if [ "$IS_UPDATE" -eq 1 ]; then
        echo "Stopping running instance(s)..."
        pkill -f "$INSTALL_PATH" 2>/dev/null || true
        sleep 1
    fi
    
    echo "Installing to $INSTALL_PATH..."
    cp "$TEMP_BIN" "$INSTALL_PATH"
    chmod +x "$INSTALL_PATH"
    
    echo "Running install..."
    "$INSTALL_PATH" install --in-place --show-ui
fi

echo ""
if [ "$SKIP_INSTALL" -eq 1 ]; then
    echo "SteamInputDB Buddy is already up to date."
    elif [ "$IS_UPDATE" -eq 1 ]; then
    echo "SteamInputDB Buddy updated successfully!"
else
    echo "SteamInputDB Buddy installed successfully!"
fi
echo "Binary: $INSTALL_PATH"
