#!/bin/bash

# exit when any command fails
set -e

export GIT_COMMIT=$(git rev-list -1 HEAD)
export GIT_TAG=$(git describe --tags --abbrev=0)

echo "Invoking build script with GIMME_ARCH: ${GIMME_ARCH}"
echo "GIT COMMIT: ${GIT_COMMIT}"
echo "GIT VERSION: ${GIT_TAG}"

if [[ $GIMME_OS == 'windows' ]]; then
    go build -ldflags "-X info.GitCommit=${GIT_COMMIT} -X info.Version=${GIT_TAG}"
    mv mwa.exe mwa_windows_amd64.exe
	exit 0
fi

if [[ $GIMME_ARCH == 'amd64' ]]; then
    echo "Running golang test for linux-x64"
    go test -v ./...
    go build -ldflags "-X info.GitCommit=${GIT_COMMIT} -X info.Version=${GIT_TAG}"
    mv mwa mwa_linux_amd64
fi

if [[ $GIMME_ARCH == 'arm' ]]; then
    echo "Running golang release build for linux-arm"
    go build -ldflags "-s -w -X info.GitCommit=${GIT_COMMIT} -X info.Version=${GIT_TAG}"
    mv mwa mwa_linux_arm
fi



