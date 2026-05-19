#!/usr/bin/env bash
set -euo pipefail

TAG="$1"
VERSION="${TAG#v}"
SHA_ARM64=$(shasum -a 256 artifacts/pomogoro-arm64.dmg | awk '{print $1}')
SHA_AMD64=$(shasum -a 256 artifacts/pomogoro-amd64.dmg | awk '{print $1}')
REPO_URL="https://x-access-token:${HOMEBREW_TAP_TOKEN}@github.com/bemoty/homebrew-tap.git"

git clone "$REPO_URL" /tmp/homebrew-tap

cat > /tmp/homebrew-tap/Casks/pomogoro.rb << EOF
cask "pomogoro" do
  version "${VERSION}"

  on_arm do
    url "https://github.com/bemoty/pomogoro/releases/download/${TAG}/pomogoro-arm64.dmg"
    sha256 "${SHA_ARM64}"
  end

  on_intel do
    url "https://github.com/bemoty/pomogoro/releases/download/${TAG}/pomogoro-amd64.dmg"
    sha256 "${SHA_AMD64}"
  end

  name "pomogoro"
  desc "Simple Pomodoro timer"
  homepage "https://bemoty.dev"

  app "pomogoro.app"
end
EOF

cd /tmp/homebrew-tap
git config user.email "josh@bemoty.dev"
git config user.name "Joshua Winkler"
git add Casks/pomogoro.rb
git commit -m "pomogoro ${TAG}"
git push
