#!/bin/bash

readonly BRANCH="${TRAVIS_BRANCH:-develop}"
readonly DIR="$(realpath $(dirname ${BASH_SOURCE[0]}))"
readonly GO_VERSION="1.9"
readonly GO_WORKSPACE="/usr/local/go/src/github.com/mdelapenya/lpn"
readonly RELEASE_VERSION=$(cat ./VERSION.txt)

CHANNEL="unstable"
VERSION="$RELEASE_VERSION-snapshot"

if [[ "$BRANCH" == "master" ]]; then
  channel="stable"
  VERSION="$RELEASE_VERSION"
fi

function bind_static_files() {
  go-bindata -pkg assets -o assets/license/license.go ./LICENSE.txt
  go-bindata -pkg assets -o assets/version/version.go ./VERSION.txt
  echo ">>> LICENSE and VERSION files bound into the binary sucessfully"
}

function build_binaries() {
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
}

function git_branch_name() {
  echo $(git symbolic-ref --short HEAD)
}

function git_checks() {
  dirty=$(git_dirty)
  if [[ "$dirty" == "1" ]]; then
    echo "The repository is dirty. Please clean the workspace before releasing."
    exit 1
  fi

  tracked=$(git_num_tracked_files)
  if [[ "$tracked" != "0" ]]; then
    echo "The repository has ${tracked} tracked files. Please commit (or clean) them the workspace before releasing."
    exit 1
  fi

  untracked=$(git_num_untracked_files)
  if [[ "$untracked" != "0" ]]; then
    echo "The repository has ${untracked} untracked files. Please clean the workspace before releasing."
    exit 1
  fi

  branch=$(git_branch_name)
  if [[ "$branch" != "master" || "$branch" != "develop" ]]; then
    echo "Please create a release from master or develop branch. You are actually in '$branch'."
    exit 1
  fi

  echo "0"
}

function git_dirty() {
  [[ $(git diff --shortstat 2>/dev/null | tail -n1) != "" ]] && echo "1"
}

function git_num_untracked_files() {
  expr $(git status --porcelain 2>/dev/null | grep "^??" | wc -l | xargs)
}

function git_num_tracked_files() {
  expr $(git status --porcelain 2>/dev/null | grep "^M" | wc -l | xargs)
}

function main() {
  read -p "The ${VERSION} release of lpn is going to be created. Are you sure (y/n)?" choice
  case "$choice" in
  y | Y)
    release
    ;;
  *)
    echo "Sorry about that. Maybe next time :( Remember to type y/Y if you want to release a new version."
    exit 1
    ;;
  esac
}

function publish_binaries() {
  cd wedeploy
  we login --no-browser
  we deploy -p lpn
}

function release() {
  bind_static_files

  gitChecks=$(git_checks)

  if [[ "$gitChecks" != "0" ]]; then
    echo "$gitChecks"
    exit 1
  fi

  result=$(git tag "$VERSION")

  if [[ "$result" != "" ]]; then
    echo "$result. Please bump a version editing VERSION.txt file. Existing tags are:"
    git tag
    exit 1
  else
    echo ">>> Git tag $VERSION created successfully."
  fi

  git push origin master --tags

  echo ">>> Release $VERSION pushed to Github successfully."

  build_binaries

  echo ">>> Binaries for $VERSION built successfully."

  publish_binaries

  echo ">>> Binaries for $VERSION published to WeDeploy successfully."
}

main
