#!/bin/bash

readonly CONTAINER_ID="liferay-portal-nook"
readonly DIR="${TRAVIS_BUILD_DIR:-.}"
readonly NIGHTLY_IMAGE="mdelapenya/liferay-portal-nightlies"
readonly LPN_GO_BINARY="${DIR}/lpn"
readonly RELEASE_IMAGE="mdelapenya/liferay-portal"

function test_run_nightly() {
  echo "Test: test_run_nightly"
  image="${NIGHTLY_IMAGE}"
  tag="latest"

  $LPN_GO_BINARY run nightly -t "${tag}"

  exists=$(checkContainer)

  if [[ "$exists" == "1" ]]; then
    echo "Container ${CONTAINER_ID} is not running."
    exit 1
  else
    echo "Container ${CONTAINER_ID} is running."
  fi

  $LPN_GO_BINARY rm
}

function test_run_release() {
  echo "Test: test_run_release"
  image="${RELEASE_IMAGE}"
  tag="latest"

  $LPN_GO_BINARY run release -t "${tag}"

  exists=$(checkContainer)

  if [[ "$exists" == "1" ]]; then
    echo "Container ${CONTAINER_ID} is not running."
    exit 1
  else
    echo "Container ${CONTAINER_ID} is running."
  fi

  $LPN_GO_BINARY rm
  exit ${exists}
}

main() {
  test_run_nightly
  test_run_release
}

function checkContainer() {
  exists=$(
    docker inspect ${CONTAINER_ID} &>/dev/null
    echo $?
  )

  echo ${exists}
}

main
