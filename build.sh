#!/bin/bash
set -e

APP_NAME="codeforgeai"
SRC="../codeforgeai.go"
OUT_DIR="build"

# Clear previous build outputs
rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

PLATFORMS=(
    "windows amd64"
    "windows arm64"
    "darwin amd64"
    "darwin arm64"
    "linux amd64"
    "linux arm64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    read -r GOOS GOARCH <<< "$PLATFORM"
    EXT=""
    [ "$GOOS" = "windows" ] && EXT=".exe"
    OUT_FILE="${OUT_DIR}/${APP_NAME}_${GOOS}_${GOARCH}${EXT}"
    echo "Building $OUT_FILE..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "$OUT_FILE" "$SRC"
done

echo "Builds complete. Files are in $OUT_DIR/"