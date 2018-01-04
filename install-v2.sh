#!/usr/bin/env sh

# Invoking this script:
#
# To install the latest stable version:
# curl https://deis.io/deis-cli/install-v2.sh | sh
#
# To install a specific released version ($VERSION):
# curl https://deis.io/deis-cli/install-v2.sh | sh -s $VERSION
#
# - download deis cli binary
# - making sure deis cli binary is executable
# - explain what was done
#

# install current version unless overridden by first command-line argument
VERSION=${1:-stable}

set -euf

check_platform_arch() {
  local supported="linux-amd64 darwin-amd64"

  if ! echo "${supported}" | tr ' ' '\n' | grep -q "${PLATFORM}-${ARCH}"; then
    cat <<EOF

The Deis Workflow CLI (deis) is not currently supported on ${PLATFORM}-${ARCH}.

See https://deis.com/workflow/ for more information.

EOF
  fi
}

PLATFORM="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
# https://storage.googleapis.com/workflow-cli-release/v2.0.0/deis-v2.0.0-darwin-386
DEIS_BIN_URL_BASE="https://storage.googleapis.com/workflow-cli-release"

if [ "${ARCH}" = "x86_64" ]; then
  ARCH="amd64"
fi

check_platform_arch

DEIS_CLI="deis-${VERSION}-${PLATFORM}-${ARCH}"
DEIS_CLI_PATH="${DEIS_CLI}"
if [ "${VERSION}" != 'stable' ]; then
  DEIS_CLI_PATH="${VERSION}/${DEIS_CLI_PATH}"
fi

echo "Downloading ${DEIS_CLI} From Google Cloud Storage..."
curl -fsSL -o deis "${DEIS_BIN_URL_BASE}/${DEIS_CLI_PATH}"

chmod +x deis

cat <<EOF

The Deis Workflow CLI (deis) is now available in your current directory.

To learn more about Deis Workflow, execute:

    $ ./deis --help

EOF