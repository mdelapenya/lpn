#!/bin/bash

readonly DIR="${TRAVIS_BUILD_DIR:-.}"
readonly CONTAINER_ID="liferay-portal-nook"
readonly LPN_GO_BINARY="${DIR}/lpn"

function test_stop_container_not_running() {
    echo "Test: test_stop_container_not_running"
    return=$(
        $LPN_GO_BINARY stop
    )

    if [[ "$return" == "" ]]; then
        echo "Container ${CONTAINER_ID} was not present and could not be stopped."
    else
        exit 1
    fi
}

function test_stop_container_running() {
    echo "Test: test_stop_container_running"
    $LPN_GO_BINARY run nightly
    return=$(
        $LPN_GO_BINARY stop
    )

    if [[ "$return" == "${CONTAINER_ID}" ]]; then
        echo "Container ${CONTAINER_ID} was running and could be stopped"
    else
        echo "Container ${CONTAINER_ID} was running and could not be stopped"
        exit 1
    fi
}

main() {
    test_stop_container_not_running
    test_stop_container_running
}

main