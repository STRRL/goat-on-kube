#!/usr/bin/env bash

set -exo pipefail

DIR=$(dirname $0)

COMMIT_HASH=$(bash "${DIR}"/commit-hash.sh)

cd ${DIR}/../ && \
    DOCKER_BUILDKIT=1 DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build -t goat-exporter:"${COMMIT_HASH}" \
    --build-arg GOAT_EXPORTER_REVISION="${COMMIT_HASH}" \
    -f ./Dockerfile ./

docker tag goat-exporter:"${COMMIT_HASH}" goat-exporter:latest

if [ ! -z ${IMAGE_TAG} ]; then
    docker tag goat-exporter:"${COMMIT_HASH}" goat-exporter:${IMAGE_TAG}
fi
