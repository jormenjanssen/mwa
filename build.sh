#!/bin/bash
echo "Invoking build script with TRAVIS_OS_NAME: ${TRAVIS_OS_NAME}"

if [[ $TRAVIS_OS_NAME == 'linux-x64' ]]; then
    echo "Running golang test for linux-x64"
    go test -v ./...
fi

if [[ $TRAVIS_OS_NAME == 'linux-arm' ]]; then
    echo "Running golang release build for linux-arm"
    go build -ldflags "-s -w"
fi



