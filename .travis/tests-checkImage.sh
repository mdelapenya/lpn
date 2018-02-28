#!/bin/bash

readonly DIR="${TRAVIS_BUILD_DIR:-.}"
readonly NIGHTLY_IMAGE="mdelapenya/liferay-portal-nightlies"
readonly LPN_GO_BINARY="${DIR}/lpn"
readonly RELEASE_IMAGE="mdelapenya/liferay-portal"

function test_checkImage_nightly() {
  echo "Test: test_checkImage_nightly"
  image="${NIGHTLY_IMAGE}"
  tag="latest"
  $LPN_GO_BINARY pull nightly -t "${tag}"
  exists=$(
    $LPN_GO_BINARY checkImage nightly -t "${tag}"
  )

  if [[ "$exists" == "1" ]]; then
    echo "Image ${image}:${tag} was not present"
    exit 1
  fi
}

function test_checkImage_not_existing() {
  echo "Test: test_checkImage_not_existing"
  image="${RELEASE_IMAGE}"
  tag="bar"
  $LPN_GO_BINARY pull release -t "${tag}"
  exists=$(
    $LPN_GO_BINARY checkImage release -t "${tag}"
  )

  if [[ "$exists" == "1" ]]; then
    echo "Image ${image}:${tag} was not present"
  fi
}

function test_checkImage_release() {
  echo "Test: test_checkImage_release"
  image="${RELEASE_IMAGE}"
  tag="latest"
  $LPN_GO_BINARY pull release -t "${tag}"
  exists=$(
    $LPN_GO_BINARY checkImage release -t "${tag}"
  )

  if [[ "$exists" == "1" ]]; then
    echo "Image ${image}:${tag} was not present"
    exit 1
  fi
}

main() {
  test_checkImage_nightly
  test_checkImage_not_existing
  test_checkImage_release
}

main
