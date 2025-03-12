#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script updates `pkg/i18n/translations/drycc/template.pot` and
# generates/fixes .po and .mo files.
# Usage: `update-translations.sh`.


# Opposite of ensure-temp-dir()
cleanup-temp-dir() {
  rm -rf "${WORKFLOW_CLI_TEMP}"
}

# Create a temp dir that'll be deleted at the end of this bash session.
#
# Vars set:
#   WORKFLOW_CLI_TEMP
ensure-temp-dir() {
  if [[ -z ${WORKFLOW_CLI_TEMP-} ]]; then
    WORKFLOW_CLI_TEMP=$(mktemp -d 2>/dev/null || mktemp -d -t kubernetes.XXXXXX)
    trap cleanup-temp-dir EXIT
  fi
}


WORFLOW_CLI_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

TRANSLATIONS="pkg/i18n/translations"
WORFLOW_CLI_FILES=()
WORFLOW_CLI_DEFAULT_LOCATIONS=(
  "cmd"
  "internal/parser"
)

generate_pot="false"
generate_mo="false"
fix_translations="false"

while getopts "hf:xkg" opt; do
  case ${opt} in
    h)
      echo "$0 [-f files] [-x] [-k] [-g]"
      echo " -f <file-path>: Files to process"
      echo " -x extract strings to a POT file"
      echo " -k fix .po files; deprecate translations by marking them obsolete and supply default messages"
      echo " -g sort .po files and generate .mo files"
      exit 0
      ;;
    f)
      WORFLOW_CLI_FILES+=("${OPTARG}")
      ;;
    x)
      generate_pot="true"
      ;;
    k)
      fix_translations="true"
      ;;
    g)
      generate_mo="true"
      ;;
    \?)
      echo "[-f <files>] -x -g" >&2
      exit 1
      ;;
  esac
done

if [[ ${#WORFLOW_CLI_FILES} -eq 0 ]]; then
  WORFLOW_CLI_FILES+=("${WORFLOW_CLI_DEFAULT_LOCATIONS[@]}")
fi

if ! which go-xgettext > /dev/null; then
  echo 'Can not find go-xgettext, install with:'
  echo 'go install github.com/gosexy/gettext/go-xgettext@latest'
  exit 1
fi

if ! which msgfmt > /dev/null; then
  echo 'Can not find msgfmt, install with:'
  echo 'apt-get install gettext'
  exit 1
fi

if [[ "${generate_pot}" == "true" ]]; then
  echo "Extracting strings to POT"
  # shellcheck disable=SC2046
  go-xgettext -k=i18n.T $(grep -lr "i18n.T" "${WORFLOW_CLI_FILES[@]}") > tmp.pot
  perl -pi -e 's/CHARSET/UTF-8/' tmp.pot
  perl -pi -e 's/\\(?!n"\n|")/\\\\/g' tmp.pot
  ensure-temp-dir
  if msgcat -s tmp.pot > "${WORKFLOW_CLI_TEMP}/template.pot"; then
    mv "${WORKFLOW_CLI_TEMP}/template.pot" "${TRANSLATIONS}/drycc/template.pot"
    rm tmp.pot
  else
    echo "Failed to update template.pot"
    exit 1
  fi
fi

if [[ "${fix_translations}" == "true" ]]; then
  echo "Fixing .po files"
  ensure-temp-dir
  for PO_FILE in "${TRANSLATIONS}"/drycc/*/LC_MESSAGES/cli.po; do
    TMP="${WORKFLOW_CLI_TEMP}/fix.po"
    msgen "${TRANSLATIONS}/drycc/template.pot" | \
      msgmerge --no-fuzzy-matching "${PO_FILE}" - > "${TMP}"
    mv "${TMP}" "${PO_FILE}"
  done
fi

if [[ "${generate_mo}" == "true" ]]; then
  echo "Generating .po and .mo files"
  for x in "${TRANSLATIONS}"/*/*/*/*.po; do
    msgcat -s "${x}" > tmp.po
    mv tmp.po "${x}"
    echo "generating .mo file for: ${x}"
    msgfmt "${x}" -o "$(dirname "${x}")/$(basename "${x}" .po).mo"
  done
fi
