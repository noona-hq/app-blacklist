#!/bin/bash

if [[ $# -eq 0 ]] ; then
    echo 'Missing arguement: image name'
    exit 0
fi

NAME=$1

if [[ "$MODE" != "PROD" ]]; then
    TAG_SUFFIX="-dev"
fi

TAG=$(git rev-parse --short HEAD)
IMG=${NAME}:${TAG}${TAG_SUFFIX}
LATEST=${NAME}:latest${TAG_SUFFIX}

echo "Building image $IMG"
docker build -t ${IMG} .

echo "Tagging as latest $LATEST"
docker tag ${IMG} ${LATEST}