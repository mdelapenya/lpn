#!/bin/bash

readonly DIR="$(realpath $(dirname ${BASH_SOURCE[0]}))"
readonly VERSION=$(cat ./VERSION.txt)

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
  if [[ "$branch" != "master" ]]; then
    echo "Please create a release from master branch. You are actually in '$branch'."
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

function pushToEquinox() {
  cat <<EOF >.equinox/config.yaml
app: app_dK5yVpq7ybd
signing-key: .equinox/equinox.key
token: $(cat .equinox/token)
platforms: [
  darwin_amd64,
  linux_amd64,
  windows_amd64
]
EOF

  equinox release \
    --config=".equinox/config.yaml" \
    --version="$(VERSION)" \
    --channel="stable" \
    github.com/mdelapenya/lpn

  echo ">>> Release $VERSION pushed to Equinox successfully."
}

function release() {
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

  pushToEquinox
}

main
