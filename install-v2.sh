#!/usr/bin/env sh

# Invoking this script:
#
# To install the latest stable version:
# curl https://raw.githubusercontent.com/drycc/workflow-cli/master/install-v2.sh | sh
#
# To install a specific released version ($VERSION):
# curl https://raw.githubusercontent.com/drycc/workflow-cli/master/install-v2.sh | sh -s $VERSION
#
# - download drycc cli binary
# - making sure drycc cli binary is executable
# - explain what was done
#

# install current version unless overridden by first command-line argument
VERSION=${1:-stable}

set -euf

check_platform_arch() {
  local supported="linux-amd64 darwin-amd64"

  if ! echo "${supported}" | tr ' ' '\n' | grep -q "${PLATFORM}-${ARCH}"; then
    cat <<EOF

The Drycc Workflow CLI (drycc) is not currently supported on ${PLATFORM}-${ARCH}.

See https://github.com/drycc/workflow-cli for more information.

EOF
  fi
}

PLATFORM="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
# https://storage.googleapis.com/drycc-workflow-cli-release/v2.18.0/drycc-v2.18.0-darwin-386
DRYCC_BIN_URL_BASE="https://storage.googleapis.com/drycc-workflow-cli-release"

if [ "${ARCH}" = "x86_64" ]; then
  ARCH="amd64"
fi

check_platform_arch

DRYCC_CLI="drycc-${VERSION}-${PLATFORM}-${ARCH}"
DRYCC_CLI_PATH="${DRYCC_CLI}"
if [ "${VERSION}" != 'stable' ]; then
  DRYCC_CLI_PATH="${VERSION}/${DRYCC_CLI_PATH}"
fi

echo "Downloading ${DRYCC_CLI} From Google Cloud Storage..."
echo "Downloading binary from here: ${DRYCC_BIN_URL_BASE}/${DRYCC_CLI_PATH}"
curl -fsSL -o drycc "${DRYCC_BIN_URL_BASE}/${DRYCC_CLI_PATH}"

chmod +x drycc

cat <<EOF

The Drycc Workflow CLI (drycc) is now available in your current directory.

To learn more about Drycc Workflow, execute:

    $ ./drycc --help

EOF