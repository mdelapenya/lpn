#!/bin/bash

function git_branch_name() {
  echo $(git symbolic-ref --short HEAD)
}

readonly BRANCH="${TRAVIS_BRANCH:-$(git_branch_name)}"
readonly DIR="$(realpath $(dirname ${BASH_SOURCE[0]}))"
readonly GO_VERSION="1.9"
readonly GO_WORKSPACE="/usr/local/go/src/github.com/mdelapenya/lpn"
readonly VERSION=$(cat ./VERSION.txt)

if [[ "$BRANCH" != "master" ]]; then
    echo "It's not possible to build from a branch different to master"
    exit 1
fi

function bind_static_files() {
  go-bindata -pkg assets -o assets/license/license.go ./LICENSE.txt
  go-bindata -pkg assets -o assets/version/version.go ./VERSION.txt
  echo ">>> LICENSE and VERSION files bound into the binary successfully"
}

function build_binaries() {
  ./scripts/build.sh
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
  if [[ "$branch" != "master" ]] && [[ "$branch" != "develop" ]]; then
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

function publish_website() {
  ./scripts/publish.sh
}

function release() {
  bind_static_files

  gitChecks=$(git_checks)

  if [[ "$gitChecks" != "0" ]]; then
    echo "$gitChecks"
    exit 1
  fi

  build_binaries

  echo ">>> Binaries for $VERSION built successfully."

  result=$(git tag "$VERSION")

  if [[ "$result" != "" ]]; then
    echo "$result. Please bump a version editing VERSION.txt file. Existing tags are:"
    git tag
    exit 1
  else
    echo ">>> Git tag $VERSION created successfully."
  fi

  git push origin $BRANCH --tags

  echo ">>> Release $VERSION pushed to Github successfully."

  publish_website

  echo ">>> Website for $VERSION published to WeDeploy successfully."
}

main
