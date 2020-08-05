include scripts/commands/vars.sh

export BRANCH ?= $(shell git branch | sed -n 's/^\* //p')
export COMMIT_DATE_UTC ?= $(shell TZ=UTC git show --quiet --date='format-local:%Y%m%d_%H%M%S' --format="%cd")

export DOCKER_TAG ?= $(shell echo ${BRANCH} | tr '[:upper:]' '[:lower:]' | sed 's,/,_,g')
export CHART_DIR ?= k8s/${APPNAME}

setup:
	@scripts/commands/setup.sh

build:
	@scripts/commands/build.sh

## Upload helm charts for deploying on k8s
helm-publish:
	@echo "Publishing helm package to Artifactory"
	helm lint ${CHART_DIR}
	helm package ${CHART_DIR}
	jfrog rt u "*.tgz" "helm-local/yapo/" || true

run:
	./${APPNAME}

start: build run

docker-build:
	@scripts/commands/docker-build.sh

docker-publish:
	@scripts/commands/docker-publish.sh

docker-attach:
	@scripts/commands/docker-attach.sh

docker-compose-up-prod:
	@scripts/commands/docker-compose-up.sh prod

docker-compose-down-prod:
	@scripts/commands/docker-compose-down.sh prod

docker-compose-up-dev:
	@scripts/commands/docker-compose-up.sh dev
	@scripts/commands/docker-attach.sh

docker-compose-down-dev:
	@scripts/commands/docker-compose-down.sh dev

validate:
	@scripts/commands/validate.sh

fix-format:
	@scripts/commands/fix-format.sh

tests:
	@scripts/commands/tests.sh

clean:
	rm -rf ${APPNAME} ${APPNAME}.* daemon.* reports vendor

info:
	@echo "Service: ${APPNAME}"
	@echo "Images from latest commit:"
	@echo "- ${DOCKER_IMAGE}:${DOCKER_TAG}"
	@echo "- ${DOCKER_IMAGE}:${COMMIT_DATE_UTC}"
	@echo "API Base URL: ${BASE_URL}"
	@echo "Healthcheck: curl ${BASE_URL}/api/v1/healthcheck"
