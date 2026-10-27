#!/usr/bin/env bash

# Traces the execution of jenai from the current source using dlv.

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORKDIR="$PWD"
cd "$ROOT"
dlv trace --wd "$WORKDIR" "./go/cmd/jenai" "$@"
