#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

########### CODE ##############

#Publishing is only allowed from Travis
if [[ -n "$TRAVIS" ]]; then
    echoTitle "Publishing docker image to Artifactory"
    docker login --username "${ARTIFACTORY_USER}" --password "${ARTIFACTORY_PWD}" "${DOCKER_REGISTRY}"
    docker push "${DOCKER_IMAGE}"
else
    echoError "DOCKER PUBLISHING IS ONLY ALLOWED IN TRAVIS"
fi
