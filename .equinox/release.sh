#!/bin/bash

readonly DIR="$(realpath $(dirname ${BASH_SOURCE[0]}))"
readonly EQUINOX_APP_ID=app_dK5yVpq7ybd
VERSION=$(cat ${DIR}/../VERSION.txt)
readonly BRANCH="${TRAVIS_BRANCH:-develop}"

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
  openssl aes-256-cbc -K $encrypted_2aabfcb2deac_key -iv $encrypted_2aabfcb2deac_iv -in equinox.key.enc -out equinox.key -d

  cat <<EOF >${DIR}/config.yaml
app: ${EQUINOX_APP_ID}
signing-key: ${DIR}/equinox.key
token: ${EQUINOX_TOKEN}
platforms: [
  darwin_amd64,
  linux_amd64,
  windows_amd64
]
EOF

  CHANNEL="stable"

  if [[ "$BRANCH" == "develop" ]]; then
    CHANNEL="unstable"
    VERSION="$VERSION-snapshot"
  fi

  equinox release \
    --config="${DIR}/config.yaml" \
    --version="$VERSION" \
    --channel="$CHANNEL" \
    github.com/liferay/lpn

  echo ">>> Release $VERSION pushed to Equinox successfully."
}

main
