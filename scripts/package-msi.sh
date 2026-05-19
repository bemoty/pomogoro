#!/usr/bin/env bash
set -euo pipefail

VERSION="${1#v}"
STAGE=$(mktemp -d)
trap 'rm -rf "$STAGE"' EXIT

cp artifacts/pomogoro-windows-amd64.exe "$STAGE/pomogoro.exe"
cp assets/pomogoro.ico "$STAGE/"
sed "s/@VERSION@/${VERSION}/g" windows/app.wxs > "$STAGE/app.wxs"

wixl -v --arch x64 -o artifacts/pomogoro-amd64.msi "$STAGE/app.wxs"
