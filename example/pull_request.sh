#!/bin/sh

DIRECTORY=${1}
GITURL=${2}
GITBRANCH=${3}

updateRepo() {
    local path="$1"
    local branch="$2"

    echo "Processing pull on" $path

    git -C $path fetch origin $branch
    git -C $path pull origin $branch
}

cloneRepo() {
    local path="$1"
    local url="$2"
    local branch="$3"

    echo "Processing clone" $url

    git clone $url -b $branch $path
}

echo "\n"
if [ ! -d "$DIRECTORY" ]; then
  echo "Directory not exists, Run clone repo"
  cloneRepo $DIRECTORY $GITURL $GITBRANCH
  else
  echo "Directory exists, Run pull repo"
  updateRepo $DIRECTORY $GITBRANCH
fi
echo ""