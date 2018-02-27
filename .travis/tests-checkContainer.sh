#!/bin/bash

readonly DIR="${TRAVIS_BUILD_DIR:-.}"
readonly NIGHTLY_IMAGE="mdelapenya/liferay-portal-nightlies"
readonly LPN_GO_BINARY="${DIR}/lpn"
readonly RELEASE_IMAGE="mdelapenya/liferay-portal"

function test_checkContainer_nightly() {
    echo "Test: test_checkContainer_nightly"
    image="${NIGHTLY_IMAGE}"
    tag="latest"
    $LPN_GO_BINARY run nightly -t "${tag}"
    $LPN_GO_BINARY checkContainer
    exit $?
}

function test_checkContainer_release() {
    echo "Test: test_checkContainer_release"
    image="${RELEASE_IMAGE}"
    tag="latest"
    $LPN_GO_BINARY pull release -t "${tag}"
    $LPN_GO_BINARY checkContainer
    exit $?
}

main() {
    test_checkContainer_nightly
    test_checkContainer_release
}

main