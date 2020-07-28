# the filepath to this repository, relative to $GOPATH/src
REPO_PATH := github.com/drycc/workflow-cli
DEV_ENV_IMAGE := drycc/go-dev
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}

# The latest git tag on branch
GIT_TAG ?= $(shell git describe --abbrev=0 --tags)
REVISION ?= $(shell git rev-parse --short HEAD)

REGISTRY ?= quay.io/
IMAGE_PREFIX ?= drycc
IMAGE := ${REGISTRY}${IMAGE_PREFIX}/workflow-cli-dev:${REVISION}

DIST_DIR ?= _dist

DEV_ENV_CMD := docker run --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_IMAGE}

define build-install-script
  sed "s|{{DRYCC-CLIENT-VERSION}}|${GIT_TAG}|g" "install.tmpl" > "${DIST_DIR}/$(1)/install-drycc.sh"
endef

bootstrap:
	${DEV_ENV_CMD} go mod vendor

# This is supposed to be run within a docker container
build-revision:
	${DEV_ENV_CMD} bash build.sh build-revision ${REVISION}

# This is supposed to be run within a docker container
build-tag:
	${DEV_ENV_CMD} bash build.sh build-revision ${GIT_TAG}
	@$(call build-install-script,${GIT_TAG})

build: build-tag build-revision

test-style:
	${DEV_ENV_CMD} lint

test-cover:
	${DEV_ENV_CMD} test-cover.sh

test: build-revision test-style test-cover
	${DEV_ENV_CMD} go test ./...
