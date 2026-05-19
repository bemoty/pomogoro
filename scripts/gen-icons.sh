#!/usr/bin/env bash
set -euo pipefail

ICONSET=$(mktemp -d)/pomogoro.iconset
mkdir -p "$ICONSET"

for size in 16 32 64 128 256 512; do
    convert -background none assets/stopwatch.svg -resize "${size}x${size}" "$ICONSET/icon_${size}x${size}.png"
done

cp "$ICONSET/icon_32x32.png"   "$ICONSET/icon_16x16@2x.png"
cp "$ICONSET/icon_64x64.png"   "$ICONSET/icon_32x32@2x.png"
cp "$ICONSET/icon_256x256.png" "$ICONSET/icon_128x128@2x.png"
cp "$ICONSET/icon_512x512.png" "$ICONSET/icon_256x256@2x.png"

iconutil -c icns "$ICONSET" -o assets/pomogoro.icns

convert \
    "$ICONSET/icon_16x16.png" \
    "$ICONSET/icon_32x32.png" \
    "$ICONSET/icon_64x64.png" \
    "$ICONSET/icon_256x256.png" \
    assets/pomogoro.ico

"$(go env GOPATH)/bin/goversioninfo" -o resource_windows_amd64.syso windows/versioninfo.json
