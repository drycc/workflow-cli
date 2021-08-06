# the filepath to this repository, relative to $GOPATH/src
REPO_PATH := github.com/drycc/workflow-cli
DEV_ENV_IMAGE := ${DEV_REGISTRY}/drycc/go-dev
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}

DIST_DIR ?= _dist

DEV_ENV_CMD := docker run --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_IMAGE}

define build-install-script
  sed "s|{{DRYCC-CLIENT-VERSION}}|${VERSION}|g" "install.tmpl" > "${DIST_DIR}/install-drycc.sh"
endef

bootstrap:
	${DEV_ENV_CMD} go mod vendor

# This is supposed to be run within a docker container
build:
	${DEV_ENV_CMD} scripts/build ${VERSION}
	@$(call build-install-script,${VERSION})

test-style:
	${DEV_ENV_CMD} lint

test-cover:
	${DEV_ENV_CMD} test-cover.sh

test: build test-style test-cover
	${DEV_ENV_CMD} go test ./...
