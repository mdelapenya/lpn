#!/bin/bash

readonly DIR="$(realpath $(dirname ${BASH_SOURCE[0]}))"

go-bindata -pkg assets -o assets/license/license.go ./LICENSE.txt
go-bindata -pkg assets -o assets/version/version.go ./VERSION.txt

echo ">>> LICENSE and VERSION files bound into the binary successfully"