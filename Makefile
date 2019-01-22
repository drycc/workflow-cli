# the filepath to this repository, relative to $GOPATH/src
repo_path = github.com/drycc/workflow-cli

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

define check-static-binary
  if file $(1) | egrep -q "(statically linked|Mach-O)"; then \
    echo -n ""; \
  else \
    echo "The binary file $(1) is not statically linked. Build canceled"; \
    exit 1; \
  fi
endef

build: build-test-image
	$(eval GO_LDFLAGS= -ldflags '-X ${repo_path}/version.Version=dev-${REVISION}')
	docker run --rm -e GOOS=${GOOS} -v ${CURDIR}/_dist:/out ${IMAGE} go build -a -installsuffix cgo ${GO_LDFLAGS} -o /out/drycc .
	@$(call check-static-binary,_dist/drycc)
	@echo "${GOOS} binary written to _dist/drycc"

# This is supposed to be run within a docker container
build-latest:
	$(eval GO_LDFLAGS = -ldflags '-X ${repo_path}/version.Version=${GIT_TAG}-${REVISION}')
	gox -verbose ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="${DIST_DIR}/drycc-latest-{{.OS}}-{{.Arch}}" .

# This is supposed to be run within a docker container
build-revision:
	$(eval GO_LDFLAGS = -ldflags '-X ${repo_path}/version.Version=${GIT_TAG}-${REVISION}')
	gox -verbose ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="${DIST_DIR}/${REVISION}/drycc-${REVISION}-{{.OS}}-{{.Arch}}" .

# This is supposed to be run within a docker container
build-stable:
	$(eval GO_LDFLAGS = -ldflags '-X ${repo_path}/version.Version=${GIT_TAG}')
	gox -verbose ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="${DIST_DIR}/drycc-stable-{{.OS}}-{{.Arch}}" .

# This is supposed to be run within a docker container
build-tag:
	$(eval GO_LDFLAGS = -ldflags '-X ${repo_path}/version.Version=${GIT_TAG}')
	gox -verbose ${GO_LDFLAGS} -os="${BUILD_OS}" -arch="${BUILD_ARCH}" -output="${DIST_DIR}/${GIT_TAG}/drycc-${GIT_TAG}-{{.OS}}-{{.Arch}}" .

build-all: build-latest build-revision

install:
	cp drycc $$GOPATH/bin

test-style: build-test-image
	docker run --rm ${IMAGE} lint

test: build-test-image
	docker run --rm ${IMAGE} test

build-test-image:
	docker build ${DOCKER_BUILD_FLAGS} -t ${IMAGE} .

push-test-image: build-test-image
	docker push ${IMAGE}
