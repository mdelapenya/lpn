#!/bin/bash

set -o errexit
set -o nounset

readonly NIGHTLY_IMAGE="mdelapenya/liferay-portal-nightlies:latest"
readonly RELEASE_IMAGE="mdelapenya/liferay-portal:7-ce-ga5-tomcat-hsql"
readonly RELEASE_IMAGE_LATEST="mdelapenya/liferay-portal:latest"

main() {
  cat <<EOF >.tmp_images
${NIGHTLY_IMAGE}
${RELEASE_IMAGE}
${RELEASE_IMAGE_LATEST}
EOF

  pull_images_concurrently "$(cat .tmp_images)"

  rm .tmp_images
}

_pull_image() {
  local image=$1
  docker pull $image
}

pull_images_concurrently() {
  local desired_images="$1"
  local name
  local version
  local image_name

  declare -a pidlist
  declare -a images
  declare -a fails

  echo "INFO: Starting to pull images."

  while read -r image_name; do
    printf "\n\e[1;31m Pulling: [$image_name] \e[0m\n"
    _pull_image $image_name &
    pidlist+=($!)
    images+=($image_name)
  done<<<"$desired_images"

  for index in "${!pidlist[@]}"; do
    local image=${images[$index]}
    local pid=${pidlist[$index]}

    echo "waiting on pid $pid"
    if ! wait $pid; then
      fails+=($image)
    fi
  done

  set +o nounset
  if ! [[ ${#fails[@]} -eq 0 ]]; then
    printf "\n\e[1;31m Some image pulls failed: \e[0m\n"
    for fail in ${fails[@]}; do
      printf "\e[1;31m - $fail \e[0m\n"
    done

    exit 1
  fi
}

stop_pulls() {
  printf "\e[1;31mStopping pulls\e[0m\n"

  kill -- -$(ps -o pgid=$$ | grep -o '[0-9]*')
}

trap stop_pulls INT

main