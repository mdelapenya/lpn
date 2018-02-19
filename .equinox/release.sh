#!/bin/bash

readonly DIR="$(realpath $(dirname ${BASH_SOURCE[0]}))"

cat <<EOF > $DIR/config.yaml
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
  --config="$DIR/config.yaml" \
  --version="$(cat ../VERSION.txt)" \
  --channel="stable" \
  github.com/mdelapenya/lpn