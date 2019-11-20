#!/usr/bin/env bash
export UNAMESTR = $(uname)
export GO_FILES = $(shell find . -iname '*.go' -type f | grep -v vendor | grep -v pact) # All the .go files, excluding vendor/
GENPORTOFF?=0
genport = $(shell expr ${GENPORTOFF} + \( $(shell id -u) - \( $(shell id -u) / 100 \) \* 100 \) \* 200 + 30200 + $(1))
export REPORT_ARTIFACTS=reports

# GIT variables
export GIT_BRANCH=$(shell git branch | sed -n 's/^\* //p')
export GIT_TAG=$(shell git tag -l --points-at HEAD | tr '\n' '_' | sed 's/_$$//;')
export GIT_COMMIT=$(shell git rev-parse HEAD)
export GIT_COMMIT_SHORT=$(shell git rev-parse --short HEAD)
export BUILD_CREATOR=$(shell git log --format=format:%ae | head -n 1)
export GIT_BRANCH_LOWERCASE=$(shell echo "${GIT_BRANCH}" | awk '{print tolower($0)}'| sed 's/\//_/;')

#APP variables
export APPNAME=rabbit2kafka
export APP_VERSION=0.0.1
export EXEC=./${APPNAME}
export SERVER_ROOT=${PWD}
export SERVERNAME=`hostname`
export MAIN_FILE=cmd/${APPNAME}/main.go
export LOGGER_SYSLOG_ENABLED=false
export LOGGER_STDLOG_ENABLED=true
export LOGGER_LOG_LEVEL=0

export RABBITMQ_HOST=localhost
export RABBITMQ_PORT=5672
export RABBITMQ_QUEUE=backend_event
export RABBITMQ_EXCHANGE=/

export KAFKA_HOST=localhost
export KAFKA_PORT=9093
export KAFKA_TOPIC=events_queue

#DOCKER variables
export DOCKER_REGISTRY=containers.schibsted.io
export DOCKER_IMAGE=${DOCKER_REGISTRY}/yapo/${APPNAME}
export DOCKER_BINARY=${APPNAME}.docker
export DOCKER_RABBITMQ_HOST=rabbit
export DOCKER_KAFKA_HOST=kafka

#DOCKER COMPOSE variables
export DOCKER_COMPOSE_NETWORKS=regress
