#!/usr/bin/env bash
export GO_FILES = $(shell find . -iname '*.go' -type f | grep -v vendor | grep -v pact) # All the .go files, excluding vendor/
GENPORTOFF?=0
genport = $(shell expr ${GENPORTOFF} + \( $(shell id -u) - \( $(shell id -u) / 100 \) \* 100 \) \* 200 + 30200 + $(1))
export REPORT_ARTIFACTS=reports

# GIT variables
export GIT_COMMIT=$(shell git rev-parse HEAD)
export GIT_COMMIT_DATE=$(shell TZ="America/Santiago" git show --quiet --date='format-local:%d-%m-%Y_%H:%M:%S' --format="%cd")
export BUILD_CREATOR=$(shell git log --format=format:%ae | head -n 1)

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

export RABBITMQ_HOST=172.21.10.240
export KAFKA_HOST=172.21.1.95
export KAFKA_TOPIC=dev01-events-queue

#DOCKER variables
export DOCKER_REGISTRY=containers.mpi-internal.com
export DOCKER_IMAGE=${DOCKER_REGISTRY}/yapo/${APPNAME}
export DOCKER_BINARY=${APPNAME}.docker
export DOCKER_RABBITMQ_HOST=rabbit
export DOCKER_KAFKA_HOST=kafka

#DOCKER COMPOSE variables
export DOCKER_COMPOSE_NETWORKS=regress
