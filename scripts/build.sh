#!/bin/bash

readonly BRANCH="${TRAVIS_BRANCH:-develop}"
readonly GO_VERSION="1.9"
readonly BIN_DIR="$(pwd)/bin"
readonly GO_WORKSPACE="/usr/local/go/src/github.com/mdelapenya/lpn"

CHANNEL="stable"
VERSION="$(cat ./VERSION.txt)"

if [[ "$BRANCH" == "develop" ]]; then
    CHANNEL="unstable"
    VERSION="$VERSION-snapshot"
fi

for GOOS in darwin linux windows; do
    extension=""

    if [[ "$GOOS" == "windows" ]]; then
        extension=".exe"
    fi

    for GOARCH in 386 amd64; do
        echo ">>> Building for ${GOOS}/${GOARCH}"
        docker run --rm -v "$(pwd)":${GO_WORKSPACE} -w ${GO_WORKSPACE} \
            -e GOOS=${GOOS} -e GOARCH=${GOARCH} golang:${GO_VERSION} \
            go build -v -o ${GO_WORKSPACE}/wedeploy/bin/${CHANNEL}/${VERSION}/${GOOS}/${GOARCH}/lpn${extension}
    done
done