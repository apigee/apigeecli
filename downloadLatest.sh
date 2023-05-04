#!/bin/sh
# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -e

# Determines the operating system.
OS="$(uname)"
if [ "${OS}" = "Darwin" ] ; then
  OSEXT="Darwin"
else
  OSEXT="Linux"
fi

# Determine the latest apigeecli version by version number ignoring alpha, beta, and rc versions.
if [ "${APIGEECLI_VERSION}" = "" ] ; then
  APIGEECLI_VERSION="$(curl -sL https://github.com/apigee/apigeecli/releases/latest | \
                  grep -i release | grep -o 'v[0-9].[0-9][0-9][0-9]' | tail -1)"
  APIGEECLI_VERSION="${APIGEECLI_VERSION##*/}"
fi

LOCAL_ARCH=$(uname -m)
if [ "${TARGET_ARCH}" ]; then
    LOCAL_ARCH=${TARGET_ARCH}
fi

case "${LOCAL_ARCH}" in
  x86_64|amd64|arm64)
    APIGEECLI_ARCH=x86_64
    ;;
  armv8*|aarch64*)
    APIGEECLI_ARCH=arm64
    ;;
  *)
    echo "This system's architecture, ${LOCAL_ARCH}, isn't supported"
    exit 1
    ;;
esac

if [ "${APIGEECLI_VERSION}" = "" ] ; then
  printf "Unable to get latest apigeecli version. Set APIGEECLI_VERSION env var and re-run. For example: export APIGEECLI_VERSION=v1.104"
  exit 1;
fi

# older versions of apigeecli used a file called .apigeecli. This changed to a folder in v1.108
APIGEECLI_FILE=~/.apigeecli
if [ -f "$APIGEECLI_FILE" ]; then
    rm ${APIGEECLI_FILE}
fi

# Downloads the apigeecli binary archive.
tmp=$(mktemp -d /tmp/apigeecli.XXXXXX)
NAME="apigeecli_$APIGEECLI_VERSION"

cd "$tmp" || exit
URL="https://github.com/apigee/apigeecli/releases/download/${APIGEECLI_VERSION}/apigeecli_${APIGEECLI_VERSION}_${OSEXT}_${APIGEECLI_ARCH}.zip"

download_cli() {
  printf "\nDownloading %s from %s ...\n" "$NAME" "$URL"
  if ! curl -o /dev/null -sIf "$URL"; then
    printf "\n%s is not found, please specify a valid APIGEECLI_VERSION and TARGET_ARCH\n" "$URL"
    exit 1
  fi
  curl -fsLO "$URL"
  filename="apigeecli_${APIGEECLI_VERSION}_${OSEXT}_${APIGEECLI_ARCH}.zip"
  unzip "${filename}"
  rm "${filename}"
}


download_cli

printf ""
printf "\napigeecli %s Download Complete!\n" "$APIGEECLI_VERSION"
printf "\n"
printf "apigeecli has been successfully downloaded into the %s folder on your system.\n" "$tmp"
printf "\n"

# setup apigeecli
cd "$HOME" || exit
mkdir -p "$HOME/.apigeecli/bin"
mv "${tmp}/apigeecli_${APIGEECLI_VERSION}_${OSEXT}_${APIGEECLI_ARCH}/apigeecli" "$HOME/.apigeecli/bin"
mv "${tmp}/apigeecli_${APIGEECLI_VERSION}_${OSEXT}_${APIGEECLI_ARCH}/LICENSE.txt" "$HOME/.apigeecli/LICENSE.txt"

printf "Copied apigeecli into the $HOME/.apigeecli/bin folder.\n"
chmod +x "$HOME/.apigeecli/bin/apigeecli"
rm -r "${tmp}"

# Print message
printf "\n"
printf "Added the apigeecli to your path with:"
printf "\n"
printf "  export PATH=\$PATH:\$HOME/.apigeecli/bin \n"
printf "\n"

export PATH=$PATH:$HOME/.apigeecli/bin