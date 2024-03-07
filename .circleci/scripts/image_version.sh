#!/bin/bash

if [[ "$MODE" != "PROD" ]]; then
    TAG_SUFFIX="-dev"
fi

TAG=$(git rev-parse --short HEAD)

echo ${TAG}${TAG_SUFFIX}