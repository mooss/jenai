#!/usr/bin/env bash

# Traces the execution of jenai from the current source using dlv.

SOURCE="${BASH_SOURCE[0]}"
while [ -L "$SOURCE" ]; do
  DIR="$(cd -P "$(dirname "$SOURCE")" && pwd)"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE"
done
ROOT="$(cd -P "$(dirname "$SOURCE")/.." && pwd)"

WORKDIR="$PWD"
cd "$ROOT"
dlv trace --wd "$WORKDIR" "./go/cmd/jenai" "$@"
