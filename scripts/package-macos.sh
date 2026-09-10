#!/usr/bin/env bash
set -euo pipefail

VERSION="${1#v}"
mkdir -p artifacts

for arch in arm64 amd64; do
    APP="artifacts/pomogoro.app"
    mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"

    cp "artifacts/pomogoro-darwin-${arch}" "$APP/Contents/MacOS/pomogoro"
    cp assets/pomogoro.icns "$APP/Contents/Resources/pomogoro.icns"
    sed "s/@VERSION@/${VERSION}/g" release/Info.plist.tpl > "$APP/Contents/Info.plist"

    hdiutil create \
        -volname "pomogoro" \
        -srcfolder "$APP" \
        -ov \
        -format UDZO \
        "artifacts/pomogoro-${arch}.dmg"

    rm -rf "$APP"
done
