#!/bin/bash

readonly DIR="${TRAVIS_BUILD_DIR:-.}"
readonly CONTAINER_ID="liferay-portal-nook"
readonly LPN_GO_BINARY="${DIR}/lpn"

function test_rm_container_not_running() {
  echo "Test: test_rm_container_not_running"
  return=$(
    $LPN_GO_BINARY rm
  )

  if [[ "$return" == "" ]]; then
    echo "Container ${CONTAINER_ID} was not present and could not be removed."
  else
    exit 1
  fi
}

function test_rm_container_running() {
  echo "Test: test_rm_container_running"
  $LPN_GO_BINARY run nightly
  return=$(
    $LPN_GO_BINARY rm
  )

  if [[ "$return" == "${CONTAINER_ID}" ]]; then
    echo "Container ${CONTAINER_ID} was running and could be removed"
  else
    echo "Container ${CONTAINER_ID} was running and could not be removed"
    exit 1
  fi
}

main() {
  test_rm_container_not_running
  test_rm_container_running
}

main
