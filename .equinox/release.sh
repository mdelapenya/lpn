#!/bin/bash

readonly VERSION=$(cat ./VERSION.txt)

function main() {
  installEquinox
  pushToEquinox
}

function installEquinox() {
  sudo apt-get install realpath -y
  curl -O https://bin.equinox.io/c/mBWdkfai63v/release-tool-stable-linux-amd64.zip
  unzip release-tool-stable-darwin-amd64.zip -d /usr/local/bin
}

function pushToEquinox() {
  cat <<EOF >.equinox/config.yaml
app: app_dK5yVpq7ybd
signing-key: .equinox/equinox.key
token: $(cat .equinox/token)
platforms: [
  darwin_amd64,
  linux_amd64,
  windows_amd64
]
EOF

  equinox release \
    --config=".equinox/config.yaml" \
    --version="$VERSION" \
    --channel="stable" \
    github.com/mdelapenya/lpn

  echo ">>> Release $VERSION pushed to Equinox successfully."
}

main
