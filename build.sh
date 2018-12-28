#!/bin/bash

# exit when any command fails
set -e

export GIT_COMMIT=$(git rev-list -1 HEAD)
export GIT_TAG=$(git describe --tags --abbrev=0)

echo "Invoking build script with GIMME_ARCH: ${GIMME_ARCH}"

if [[ $GIMME_OS == 'windows' ]]; then
    go build -ldflags "-X main.GitCommit=${GIT_COMMIT} main.Version=${GIT_TAG}"
	exit 0
fi

if [[ $GIMME_ARCH == 'amd64' ]]; then
    echo "Running golang test for linux-x64"
    go test -v ./...
    go build -ldflags "-X main.GitCommit=${GIT_COMMIT} main.Version=${GIT_TAG}"
    mv mwa mwa64
fi

if [[ $GIMME_ARCH == 'arm' ]]; then
    echo "Running golang release build for linux-arm"
    go build -ldflags "-s -w -X main.GitCommit=${GIT_COMMIT} main.Version=${GIT_TAG}"
fi



