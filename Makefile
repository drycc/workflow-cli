# the filepath to this repository, relative to $GOPATH/src
REPO_PATH := github.com/drycc/workflow-cli
DEV_ENV_IMAGE := quay.io/drycc/go-dev:v0.22.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}

HOST_OS := $(shell uname)
ifeq ($(HOST_OS),Darwin)
	GOOS=darwin
else
	GOOS=linux
endif

# The latest git tag on branch
GIT_TAG ?= $(shell git describe --abbrev=0 --tags)
REVISION ?= $(shell git rev-parse --short HEAD)

REGISTRY ?= quay.io/
IMAGE_PREFIX ?= drycc
IMAGE := ${REGISTRY}${IMAGE_PREFIX}/workflow-cli-dev:${REVISION}

BUILD_OS ?=linux darwin windows
BUILD_ARCH ?=amd64 386

DIST_DIR ?= _dist

DEV_ENV_CMD := docker run --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_IMAGE}

define build-install-script
  sed "s|{{DRYCC-CLIENT-VERSION}}|${GIT_TAG}|g" "install.tmpl" > "${DIST_DIR}/$(1)/install-drycc.sh"
endef

bootstrap:
	${DEV_ENV_CMD} dep ensure

# This is supposed to be run within a docker container
build-revision:
	${DEV_ENV_CMD} gox -verbose ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="${DIST_DIR}/${REVISION}/drycc-${REVISION}-{{.OS}}-{{.Arch}}" .

# This is supposed to be run within a docker container
build-tag:
	${DEV_ENV_CMD} gox -verbose ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="${DIST_DIR}/${GIT_TAG}/drycc-${GIT_TAG}-{{.OS}}-{{.Arch}}" .
	@$(call build-install-script,${GIT_TAG})

build: build-tag build-revision

test-style: build-test-image
	${DEV_ENV_CMD} lint

test:
	${DEV_ENV_CMD} test
