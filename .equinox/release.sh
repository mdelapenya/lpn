#!/bin/bash

readonly DIR="$(realpath $(dirname ${BASH_SOURCE[0]}))"
readonly VERSION=$(cat ${DIR}/../VERSION.txt)

function main() {
  installEquinox
  pushToEquinox
}

function installEquinox() {
  sudo apt-get install realpath -y
  curl -O https://bin.equinox.io/c/mBWdkfai63v/release-tool-stable-linux-amd64.zip
  sudo unzip release-tool-stable-linux-amd64.zip -d /usr/local/bin
}

function pushToEquinox() {
  cat <<EOF >${DIR}/config.yaml
app: app_dK5yVpq7ybd
signing-key: ${DIR}/equinox.key
token: $(cat ${DIR}/token)
platforms: [
  darwin_amd64,
  linux_amd64,
  windows_amd64
]
EOF

  equinox release \
    --config="${DIR}/config.yaml" \
    --version="$VERSION" \
    --channel="stable" \
    github.com/mdelapenya/lpn

  echo ">>> Release $VERSION pushed to Equinox successfully."
}

main
