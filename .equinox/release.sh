#!/bin/bash

readonly VERSION=$(cat ./VERSION.txt)

function main() {
  pushToEquinox
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
