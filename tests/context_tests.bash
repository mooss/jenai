#!/bin/bash
set -euo pipefail

#############
# Functions #
function runp() {
  go run ./go/cmd/jenai p review_previous_commit --dry-run "$@"
}

function runed() {
  go run ./go/cmd/jenai ed --dry-run "$@"
}

#########
# Setup #
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

mkdir -p "$TMP_DIR/testdir"
FILE1="$TMP_DIR/testdir/file1.txt"
FILE2="$TMP_DIR/testdir/file2.txt"
NESTED="$TMP_DIR/testdir/nested.txt"
echo "file1 content" > "$FILE1"
echo "file2 content" > "FILE2"
echo "nested content" > "$NESTED"

#########
# Tests #
# Single file inclusion.
echo -e "\nTEST: Single file context"
output=$(runp --read "$FILE1")
echo "$output" | grep -q "### Start of .*$FILE1" || {
    echo "Single file context failed"
    echo "$output"
    exit 1
}

# Multiple files.
echo -e "\nTEST: Multiple file context"
output=$(runp --read "$FILE1" "FILE2")
echo "$output" | grep -q "file1 content" || {
    echo "File1 content missing"
    exit 1
}
echo "$output" | grep -q "file2 content" || {
    echo "File2 content missing"
    exit 1
}

# Directory inclusion.
echo -e "\nTEST: Directory context"
output=$(runp --read-dir "$TMP_DIR/testdir")
echo "$output" | grep -q "nested content" || {
    echo "Directory context missing nested file"
    exit 1
}

# Context placement (above/below).
echo -e "\nTEST: Context placement"
output_above=$(runp --context-above --read-dir "$TMP_DIR/testdir")
output_below=$(runp --read "$FILE1")

if ! echo "$output_above" | head -n1 | grep -q "^## Additional context"; then
    echo "Context above failed"
    exit 1
fi

if ! echo "$output_below" | tail -n1 | grep -q "^> End of .*$FILE1"; then
    echo "Context below failed"
    exit 1
fi

# Creatable files section.
echo -e "\nTEST: Creatable files context"
CREATABLE="$TMP_DIR/testdir/newfile.txt"
rm -f "$CREATABLE"
output=$(runed arg --write "$CREATABLE")
echo "$output" | grep -q "## Creatable files" || {
    echo "Creatable files section missing"
    echo "$output"
    exit 1
}
echo "$output" | grep -q "$CREATABLE" || {
    echo "Creatable file path missing"
    echo "$output"
    exit 1
}

# Writable files section.
echo -e "\nTEST: Writable files context"
output=$(runed arg --write "$FILE1" --force)
echo "$output" | grep -q "## Editable files" || {
    echo "Editable files section missing"
    echo "$output"
    exit 1
}
echo "$output" | grep -q "file1 content" || {
    echo "Writable file content missing"
    echo "$output"
    exit 1
}

echo -e "\nAll context tests passed!"
