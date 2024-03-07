#!/bin/bash

if [[ $# != 2 ]] ; then
    echo 'Need 2 arguments: chart-version docker-image-version'
    exit 0
fi

if [[ "$MODE" == "PROD" ]]
then
    CHART_REGISTRY="helm"
else
    CHART_REGISTRY="helm-dev"
fi

CHART_VERSION=$1
IMAGE_VERSION=$2
CHART_NAME=$CIRCLE_PROJECT_REPONAME


helm chart push europe-docker.pkg.dev/timatal1-198022/${CHART_REGISTRY}/${CHART_NAME}:${CHART_VERSION}