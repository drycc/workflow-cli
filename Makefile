# the filepath to this repository, relative to $GOPATH/src
VERSION ?= canary
REPO_PATH := github.com/drycc/workflow-cli
DEV_ENV_IMAGE := ${DEV_REGISTRY}/drycc/go-dev
DEV_ENV_WORK_DIR := /opt/drycc/go/src/${REPO_PATH}

DIST_DIR ?= _dist

DEV_ENV_CMD := podman run --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} -e CODECOV_TOKEN=${CODECOV_TOKEN} ${DEV_ENV_IMAGE}

bootstrap:
	${DEV_ENV_CMD} go mod vendor

# This is supposed to be run within a container
build:
	${DEV_ENV_CMD} scripts/build ${VERSION}

test-style:
	${DEV_ENV_CMD} lint

test-cover:
	${DEV_ENV_CMD} test-cover.sh

test: build test-style test-cover
	${DEV_ENV_CMD} go test -race -cover -coverprofile=coverage.txt -covermode=atomic ./...
