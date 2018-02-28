#!/bin/bash

readonly DIR="${TRAVIS_BUILD_DIR:-.}"
readonly NIGHTLY_IMAGE="mdelapenya/liferay-portal-nightlies"
readonly LPN_GO_BINARY="${DIR}/lpn"
readonly RELEASE_IMAGE="mdelapenya/liferay-portal"

function test_pull_nightly() {
  echo "Test: test_pull_nightly"
  image="${NIGHTLY_IMAGE}"
  tag="latest"
  $LPN_GO_BINARY pull nightly -t "${tag}"

  checkImage ${image} ${tag}
}

function test_pull_release() {
  echo "Test: test_pull_release"
  image="${RELEASE_IMAGE}"
  tag="latest"
  $LPN_GO_BINARY pull release -t "${tag}"

  checkImage ${image} ${tag}
}

main() {
  test_pull_nightly
  test_pull_release
}

function checkImage() {
  image=${1}
  tag=${2}

  exists=$(
    docker inspect ${image}:${tag} &>/dev/null
    echo $?
  )

  if [[ "$exists" == "1" ]]; then
    echo "Image ${image}:${tag} is not present after pull."

    exit 1
  else
    echo "Image ${image}:${tag} is present after pull"
  fi
}

main
