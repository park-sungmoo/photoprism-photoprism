#!/usr/bin/env bash

set -e

# see https://docs.docker.com/develop/develop-images/build_enhancements/#to-enable-buildkit-builds
export DOCKER_BUILDKIT=1

if [[ -z $1 ]] && [[ -z $2 ]]; then
    echo "Please provide a container image name and version" 1>&2
    exit 1
elif [[ $1 ]] && [[ -z $2 ]]; then
    DOCKER_TAG=$(date -u +%Y%m%d)
    echo "Building 'photoprism/$1:preview'...";
    docker build --no-cache --build-arg BUILD_TAG="${DOCKER_TAG}" -t photoprism/$1:preview -f docker/${1/-//}/Dockerfile .
    echo "Done"
else
    echo "Building 'photoprism/$1:$2'...";
    docker build --no-cache --build-arg BUILD_TAG=$2 -t photoprism/$1:latest -t photoprism/$1:$2 -f docker/${1/-//}/Dockerfile .
    echo "Done"
fi
