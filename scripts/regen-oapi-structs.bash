#!/usr/bin/env bash

set -euo pipefail

mkdir -pv .workdir
cd .workdir
ROOT=..
RESOLVED=openrouter-resolved.openapi.json
RAWSPECS=openrouter.openapi.json

function curlit() {
    local source="$1"; shift
    local dest="$1";  shift
    curl --fail --show-error --location\
         --write-out '%{http_code}'\
         --output "$dest"\
         "$source" "$@"
}

function download-specs() {
    local -r url=https://openrouter.ai/openapi.json
    local -r tmp=new.json

    if [[ -e "$RAWSPECS" ]]; then
        status=$(curlit "$url" "$tmp" --time-cond "$RAWSPECS") || exit 1
    else
        status=$(curlit "$url" "$tmp") || exit 1
    fi

    if [[ "$status" = 200 ]]; then
        if [[ -e "$RAWSPECS" ]]; then
            echo Updated $RAWSPECS:
            diff -u "$RAWSPECS" "$tmp" || :
        else
            echo Downloaded $RAWSPECS for the first time.
        fi
        mv "$tmp" "$RAWSPECS"
    fi
}

function ensure-resolved() {
    download-specs
    go run $ROOT/scripts/resolve-openapi-refs.go -ref-siblings merge $RAWSPECS > $RESOLVED
}

function path() {
    jq ".paths[\"$1\"]" "$RESOLVED"
}

# gen-or generates the types for openrouter.
function gen-or() {
    local -r oapi_path="$1"; shift
    local -r dest="$1"; shift
    path "$oapi_path"\
        | python $ROOT/scripts/vibe/oapi2go.py --package generated "$@"\
        | gofmt > $ROOT/go/llm/openrouter/generated/"$dest"
}

ensure-resolved
gen-or '/model/{author}/{slug}' model.go --keep 200
gen-or '/models/{author}/{slug}/endpoints' models-endpoints.go --keep 200
