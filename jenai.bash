#!/usr/bin/env bash

HERE=$(dirname $(readlink -m "$0"))
go run "$HERE"/go/cmd/jenai "$@"
