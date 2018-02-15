#!/bin/bash

equinox release \
  --version="0.1.0" \
  --platforms="darwin_amd64 linux_amd64" \
  --signing-key=.equinox/equinox.key \
  --app="app_dK5yVpq7ybd" \
  --token="$(cat .equinox/token)" \
  github.com/mdelapenya/lpn