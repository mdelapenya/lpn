#!/bin/bash

go-bindata -pkg assets -o assets/license/license.go ./LICENSE.txt
go-bindata -pkg assets -o assets/version/version.go ./VERSION.txt

echo ">>> LICENSE and VERSION files bound into the binary successfully"