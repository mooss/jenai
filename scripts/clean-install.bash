#!/usr/bin/env bash

set -euo pipefail

COMMIT=${1:-$(git rev-parse HEAD)}
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
git clone -q . "$WORK"
cd "$WORK"
git checkout -q "$COMMIT"
echo "Building commit $(git rev-parse $COMMIT)"
BIN_DIR=$(go env GOBIN)
if [[ -z "$BIN_DIR" ]]; then
  BIN_DIR="$(go env GOPATH)/bin"
fi
go install ./go/cmd/jenai
echo "Installed to $BIN_DIR/jenai"
