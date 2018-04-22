#!/bin/bash

if [[ "${MONOREPO_SERVICE}" = "" ]]; then 
    ./GoTest.sh ${CIRCLE_SHA1}
    git-test commit ${CIRCLE_SHA1}
else 
    echo "testing service ${MONOREPO_SERVICE}"
    git-test service ${MONOREPO_SERVICE}
fi