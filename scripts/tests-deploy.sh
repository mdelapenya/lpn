#!/bin/bash

readonly CONTAINER_ID="liferay-portal-nook"
readonly DIR="${TRAVIS_BUILD_DIR:-.}"
readonly NIGHTLY_IMAGE="mdelapenya/liferay-portal-nightlies"
readonly LPN_GO_BINARY="${DIR}/lpn"
readonly RELEASE_IMAGE="mdelapenya/liferay-portal"
readonly DEPLOY_FILE_A="a.txt"
readonly DEPLOY_FILE_B="b.txt"

function test_deploy_nightly() {
  echo "Test: test_deploy_nightly"
  image="${NIGHTLY_IMAGE}"
  tag="latest"

  $LPN_GO_BINARY run nightly -t "${tag}"

  $LPN_GO_BINARY deploy nightly -f ${DIR}/scripts/resources/a.txt

  exists=$(
    docker exec ${CONTAINER_ID} ls -l /liferay/deploy | grep "${DEPLOY_FILE_A}" | wc -l | xargs
  )

  if [[ "$exists" == "1" ]]; then
    echo "File ${DEPLOY_FILE_A} has been deployed."
  else
    echo "File ${DEPLOY_FILE_A} has not been deployed."
    $LPN_GO_BINARY rm
    exit 1
  fi

  $LPN_GO_BINARY rm
}

function test_deploy_nightly_multiple_files() {
  echo "Test: test_deploy_nightly_multiple_files"
  image="${NIGHTLY_IMAGE}"
  tag="latest"

  $LPN_GO_BINARY run nightly -t "${tag}"

  $LPN_GO_BINARY deploy nightly -f ${DIR}/scripts/resources/a.txt,${DIR}/scripts/resources/b.txt

  # As file a.txt has been checked before, here we only check the second file
  exists=$(
    docker exec ${CONTAINER_ID} ls -l /liferay/deploy | grep "${DEPLOY_FILE_B}" | wc -l | xargs
  )

  if [[ "$exists" == "1" ]]; then
    echo "File ${DEPLOY_FILE_B} has been deployed."
  else
    echo "File ${DEPLOY_FILE_B} has not been deployed."
    $LPN_GO_BINARY rm
    exit 1
  fi

  $LPN_GO_BINARY rm
}

function test_deploy_release() {
  echo "Test: test_deploy_release"
  image="${RELEASE_IMAGE}"
  tag="7-ce-ga5-tomcat-hsql"

  $LPN_GO_BINARY run release -t "${tag}"

  create_deploy_folder ${tag}

  $LPN_GO_BINARY deploy release -f ${DIR}/scripts/resources/a.txt

  exists=$(
    docker exec ${CONTAINER_ID} ls -l /usr/local/${tag}/deploy | grep "${DEPLOY_FILE_A}" | wc -l | xargs
  )

  if [[ "$exists" == "1" ]]; then
    echo "File ${DEPLOY_FILE_A} has been deployed."
  else
    echo "File ${DEPLOY_FILE_A} has not been deployed."
    $LPN_GO_BINARY rm
    exit 1
  fi

  $LPN_GO_BINARY rm
}

function test_deploy_release_multiple_files() {
  echo "Test: test_deploy_release_multiple_files"
  image="${RELEASE_IMAGE}"
  tag="7-ce-ga5-tomcat-hsql"

  $LPN_GO_BINARY run release -t "${tag}"

  create_deploy_folder ${tag}

  $LPN_GO_BINARY deploy release -f ${DIR}/scripts/resources/a.txt,${DIR}/scripts/resources/b.txt

  exists=$(
    docker exec ${CONTAINER_ID} ls -l /usr/local/${tag}/deploy | grep "${DEPLOY_FILE_B}" | wc -l | xargs
  )

  if [[ "$exists" == "1" ]]; then
    echo "File ${DEPLOY_FILE_B} has been deployed."
  else
    echo "File ${DEPLOY_FILE_B} has not been deployed."
    $LPN_GO_BINARY rm
    exit 1
  fi

  $LPN_GO_BINARY rm
}

main() {
  test_deploy_nightly
  test_deploy_nightly_multiple_files

  test_deploy_release
  test_deploy_release_multiple_files
}

function create_deploy_folder() {
  tag=${1}
  deploy="/usr/local/${tag}/deploy"

  echo "Creating ${deploy} folder to avoid waiting for Liferay Portal to create it"
  docker exec ${CONTAINER_ID} mkdir -p ${deploy}
  docker exec ${CONTAINER_ID} ls -l ${deploy}
}

main
