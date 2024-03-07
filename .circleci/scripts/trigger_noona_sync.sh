#!/bin/bash

set -e

if [[ "$MODE" == "PROD" ]]
then
    BRANCH="main"
else
    BRANCH="dev"
fi

echo "Triggering sync on $BRANCH in noona-deployment"

RESP=$(curl -s -o /dev/null -w "%{http_code}" --location --request POST "https://circleci.com/api/v2/project/github/noona-hq/noona-deployment/pipeline" \
    --header "Content-Type: application/json" \
    -d "{\"branch\": \"${BRANCH}\"}" \
    -u "${CIRCLECI_TOKEN}:")

if [[ "$RESP" != 201 ]]
then
    echo "Did not get status code 201 from CircleCI"
    echo "Got $RESP"
    exit 1
fi

echo "Success"
