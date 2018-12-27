#!/bin/bash
echo "Invoking build script with GIMME_ARCH: ${GIMME_ARCH}"

if [[ $GIMME_ARCH == 'amd64' ]]; then
    echo "Running golang test for linux-x64"
    go test -v ./...
fi

if [[ $GIMME_ARCH == 'arm' ]]; then
    echo "Running golang release build for linux-arm"
    go build -ldflags "-s -w"
fi



