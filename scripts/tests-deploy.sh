#!/bin/bash

readonly CONTAINER_ID="liferay-portal-nook"
readonly DIR="${TRAVIS_BUILD_DIR:-.}"
readonly NIGHTLY_IMAGE="mdelapenya/liferay-portal-nightlies"
readonly LPN_GO_BINARY="${DIR}/lpn"
readonly RELEASE_IMAGE="mdelapenya/liferay-portal"
readonly RELEASE_HOME="liferay-ce-portal-7.0-ga5"
readonly DEPLOY_FILE_A="a.txt"
readonly DEPLOY_FILE_B="b.txt"

function test_deploy_nightly() {
  echo "Test: test_deploy_nightly"
  image="${NIGHTLY_IMAGE}"
  tag="latest"

  docker run -d  --name ${CONTAINER_ID} ${image}:${tag}

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

function test_deploy_nightly_directory() {
  echo "Test: test_deploy_nightly_directory"
  image="${NIGHTLY_IMAGE}"
  tag="latest"

  docker run -d --name ${CONTAINER_ID} ${image}:${tag}

  $LPN_GO_BINARY deploy nightly -d ${DIR}/scripts/resources

  files=( $(find ${DIR}/scripts/resources -maxdepth 1 -not -type d -and -not -name '.*' -exec basename {} \;) )

  for file in "${files[@]}"
  do
    exists=$(
      docker exec ${CONTAINER_ID} ls -l /liferay/deploy | grep "${file}" | wc -l | xargs
    )

    if [[ "$exists" != "1" ]]; then
      echo "File ${file} has not been deployed."
      $LPN_GO_BINARY rm
      exit 1
    fi
  done

  directories=( $(find ${DIR}/scripts/resources -mindepth 1 -maxdepth 1 -type d -exec basename {} \;) )

  for directory in "${directories[@]}"
  do
    exists=$(
      docker exec ${CONTAINER_ID} ls -l /liferay/deploy | grep "${directory}" | wc -l | xargs
    )

    if [[ "$exists" != "0" ]]; then
      echo "Directory ${directory} has been deployed, which was wrong."
      $LPN_GO_BINARY rm
      exit 1
    else
      echo "Directory ${directory} skipped. Cool!."
    fi
  done

  $LPN_GO_BINARY rm
}

function test_deploy_nightly_multiple_files() {
  echo "Test: test_deploy_nightly_multiple_files"
  image="${NIGHTLY_IMAGE}"
  tag="latest"

  docker run -d  --name ${CONTAINER_ID} ${image}:${tag}

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

function test_deploy_nightly_no_flag_returns_error() {
  echo "Test: test_deploy_nightly_no_flag_returns_error"
  $LPN_GO_BINARY deploy nightly

  if [[ "$?" == "1" ]]; then
    echo "Deploy nightly invoked with no commands successfully."
  else
    echo "Deploy nightly invoked with no commands did not raise an error."
    exit 1
  fi
}

function test_deploy_release() {
  echo "Test: test_deploy_release"
  image="${RELEASE_IMAGE}"
  tag="7-ce-ga5-tomcat-hsql"

  docker run -d  --name ${CONTAINER_ID} ${image}:${tag}

  create_deploy_folder ${tag}

  $LPN_GO_BINARY deploy release -f ${DIR}/scripts/resources/a.txt

  exists=$(
    docker exec ${CONTAINER_ID} ls -l /usr/local/${RELEASE_HOME}/deploy | grep "${DEPLOY_FILE_A}" | wc -l | xargs
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

function test_deploy_release_directory() {
  echo "Test: test_deploy_release_directory"
  image="${RELEASE_IMAGE}"
  tag="7-ce-ga5-tomcat-hsql"

  docker run -d --name ${CONTAINER_ID} ${image}:${tag}

  create_deploy_folder ${tag}

  $LPN_GO_BINARY deploy release -d ${DIR}/scripts/resources

  files=( $(find ${DIR}/scripts/resources -maxdepth 1 -not -type d -and -not -name '.*' -exec basename {} \;) )

  for file in "${files[@]}"
  do
    exists=$(
      docker exec ${CONTAINER_ID} ls -l /usr/local/${RELEASE_HOME}/deploy | grep "${file}" | wc -l | xargs
    )

    if [[ "$exists" != "1" ]]; then
      echo "File ${file} has not been deployed."
      $LPN_GO_BINARY rm
      exit 1
    fi
  done

  directories=( $(find ${DIR}/scripts/resources -mindepth 1 -maxdepth 1 -type d -exec basename {} \;) )

  for directory in "${directories[@]}"
  do
    exists=$(
      docker exec ${CONTAINER_ID} ls -l /usr/local/${RELEASE_HOME}/deploy | grep "${directory}" | wc -l | xargs
    )

    if [[ "$exists" != "0" ]]; then
      echo "Directory ${directory} has been deployed, which was wrong."
      $LPN_GO_BINARY rm
      exit 1
    else
      echo "Directory ${directory} skipped. Cool!."
    fi
  done

  $LPN_GO_BINARY rm
}

function test_deploy_release_multiple_files() {
  echo "Test: test_deploy_release_multiple_files"
  image="${RELEASE_IMAGE}"
  tag="7-ce-ga5-tomcat-hsql"

  docker run -d  --name ${CONTAINER_ID} ${image}:${tag}

  create_deploy_folder ${tag}

  $LPN_GO_BINARY deploy release -f ${DIR}/scripts/resources/a.txt,${DIR}/scripts/resources/b.txt

  exists=$(
    docker exec ${CONTAINER_ID} ls -l /usr/local/${RELEASE_HOME}/deploy | grep "${DEPLOY_FILE_B}" | wc -l | xargs
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

function test_deploy_release_no_flag_returns_error() {
  echo "Test: test_deploy_release_no_flag_returns_error"
  $LPN_GO_BINARY deploy release

  if [[ "$?" == "1" ]]; then
    echo "Deploy release invoked with no commands successfully."
  else
    echo "Deploy release invoked with no commands did not raise an error."
    exit 1
  fi
}

main() {
  test_deploy_nightly
  test_deploy_nightly_directory
  test_deploy_nightly_multiple_files
  test_deploy_nightly_no_flag_returns_error

  test_deploy_release
  test_deploy_release_directory
  test_deploy_release_multiple_files
  test_deploy_release_no_flag_returns_error
}

function create_deploy_folder() {
  tag=${1}
  deploy="/usr/local/${RELEASE_HOME}/deploy"

  echo "Creating ${deploy} folder to avoid waiting for Liferay Portal to create it"
  docker exec ${CONTAINER_ID} mkdir -p ${deploy}
}

main
