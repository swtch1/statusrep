#!/usr/bin/env bash

set -euo pipefail
# change this so that arguments are parsed in a more intuitive way
# search "Unofficial Bash strict mode"
IFS=$'\n\t'

export GO111MODULE=on

# make sure we're in the correct directory
cd "$(dirname "$0")"

APP="$(basename "$(pwd)")"

if [[ ! -d ./vendor ]];then
  echo 'vendoring required: use `go mod vendor` to push dependencies to vendor directory'
fi

VERSION=$(head -n 1 VERSION)
if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+.*$ ]]; then
  echo "version '$VERSION' does not match expected format, inspect VERSION file"
  exit 1
fi
VER="$VERSION-$(git log --date=raw | grep Date | head -1 | awk '{ print $2 }')"

if ! go test -v ./... -mod=vendor; then
  echo "go test failed, build process aborted"
fi

rm -rf ./bin

if ! go build -ldflags "-X main.buildVersion=${VER}" -v -mod=vendor -o "./bin/${APP}" .; then
  exit 1
fi

echo "binary written to ./bin"
