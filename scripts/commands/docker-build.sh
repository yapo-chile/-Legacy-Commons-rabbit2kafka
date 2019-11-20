#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

########### DYNAMIC VARS ###############

#In case we are in travis, docker tag will be "branch_name-20180101-1200". In case of master branch, branch_name is blank.
#In case of local build (not in travis) tag will be "local".

if [[ -n "$TRAVIS" ]]; then
    if [ "${GIT_BRANCH}" != "master" ]; then
        DOCKER_TAG=$(echo ${GIT_BRANCH} | tr '[:upper:]' '[:lower:]' | sed 's,/,_,g')
    else
        DOCKER_TAG=$(TZ=UTC git show --quiet --date='format-local:%Y%m%d_%H%M%S' --format="%cd")
    fi
else
    DOCKER_TAG=local
fi

#In case we are in travis, we will use cached docker environment.
if [[ -n "$TRAVIS" ]]; then
    DOCKER_COMMAND=container_cache
else
    DOCKER_COMMAND=docker
fi

########### CODE ##############
#Build code again now for docker platform
echoHeader "Building code for docker platform"
set -e

rm -f ${DOCKER_BINARY}
GOOS=linux GOARCH=386 go build -v -o ${DOCKER_BINARY} ./${MAIN_FILE}

set +e

echoTitle "Starting Docker Engine"
if [[ $OSTYPE == "darwin"* ]]; then
    echoTitle "Starting Mac OSX Docker Daemon"
    $DIR/docker-start-macosx.sh
elif [[ "$OSTYPE" == "linux-gnu" ]]; then
    echoTitle "Starting Linux Docker Daemon"
    sudo start-docker-daemon
else
    echoError "Platform not supported"
fi

if [[ "$BUILD_BRANCH" != "" ]]; then
    export GIT_BRANCH=${BUILD_BRANCH}
fi

echoTitle "Building docker image for ${DOCKER_IMAGE}"
echo "GIT BRANCH: ${GIT_BRANCH}"
echo "GIT TAG: ${GIT_TAG}"
echo "GIT COMMIT: ${GIT_COMMIT}"
echo "GIT COMMIT SHORT: ${GIT_COMMIT_SHORT}"
echo "BUILD CREATOR: ${BUILD_CREATOR}"
echo "BUILD NAME: ${DOCKER_IMAGE}:${GIT_COMMIT_SHORT}"


DOCKER_ARGS=" \
	-t ${DOCKER_IMAGE}:${DOCKER_TAG} \
	--build-arg GIT_BRANCH="$GIT_BRANCH" \
    --build-arg GIT_COMMIT="$GIT_COMMIT" \
    --build-arg BUILD_CREATOR="$BUILD_CREATOR" \
    --build-arg VERSION="$APP_VERSION" \
    --build-arg APPNAME="$APPNAME" \
    --build-arg BINARY="${DOCKER_BINARY}" \
    -f docker/dockerfile \
    ."

echo "args: ${DOCKER_ARGS}"
set -x
${DOCKER_COMMAND} build ${DOCKER_ARGS}
set +x

echoTitle "Build done"
