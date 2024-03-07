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

sed -i "s/version: .*/version: ${CHART_VERSION}/g" helm/${CHART_NAME}/Chart.yaml
sed -i "s/tag: .*/tag: \"${IMAGE_VERSION}\"/g" helm/${CHART_NAME}/values.yaml

helm chart save helm/${CHART_NAME} europe-docker.pkg.dev/timatal1-198022/${CHART_REGISTRY}/${CHART_NAME}