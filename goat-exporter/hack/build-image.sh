#!/usr/bin/env bash

set -exo pipefail

DIR=$(dirname $0)

COMMIT_HASH=$(bash "${DIR}"/commit-hash.sh)

IMAGE_NAME="ghcr.io/strrl/goat-on-kube/goat-exporter"

cd ${DIR}/../ && \
    DOCKER_BUILDKIT=1 DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build -t ${IMAGE_NAME}:"${COMMIT_HASH}" \
    --build-arg GOAT_EXPORTER_REVISION="${COMMIT_HASH}" \
    -f ./Dockerfile ./

docker tag ${IMAGE_NAME}:"${COMMIT_HASH}" ${IMAGE_NAME}:latest

if [ ! -z ${IMAGE_TAG} ]; then
    docker tag ${IMAGE_NAME}:"${COMMIT_HASH}" ${IMAGE_NAME}:${IMAGE_TAG}
fi
