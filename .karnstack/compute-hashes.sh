#!/usr/bin/env bash
# Emit SHA-256 hashes for every karnstack-canonical file as JSON:
#
#   {
#     ".github/workflows/verify-stages.yml": "<sha256hex>",
#     ".karnstack/compute-hashes.sh":        "<sha256hex>",
#     ".karnstack/parse-stages.sh":          "<sha256hex>",
#     ".mise.toml":                          "<sha256hex>",
#     "bloom/stage01_bit_array_test.go":     "<sha256hex>",
#     ...
#   }
#
# karnstack's /api/v1/byox/verify endpoint compares this map against the
# canonical set recorded for this template version. Any mismatch rejects
# the verification with reason `files_modified`.
#
# Files included:
#   - the workflow file itself + the two helper scripts (no smuggling
#     in an alternate verifier)
#   - .mise.toml (no rewriting the test task to no-op)
#   - every bloom/stage*_test.go (no editing the tests to pass trivially)
#
# Files NOT included: bloom/bloom.go (your implementation), README.md,
# LICENSE, go.mod, go.sum.

set -euo pipefail

if command -v sha256sum >/dev/null 2>&1; then
  hash_file() { sha256sum "$1" | awk '{print $1}'; }
elif command -v shasum >/dev/null 2>&1; then
  hash_file() { shasum -a 256 "$1" | awk '{print $1}'; }
else
  echo "neither sha256sum nor shasum found" >&2
  exit 1
fi

declare -a FILES=(
  ".github/workflows/verify-stages.yml"
  ".karnstack/compute-hashes.sh"
  ".karnstack/parse-stages.sh"
  ".mise.toml"
)

shopt -s nullglob
for f in bloom/stage*_test.go; do
  FILES+=("$f")
done
shopt -u nullglob

# Build a sorted JSON object so the output is stable across runs.
JSON="{}"
for f in "${FILES[@]}"; do
  if [ ! -f "$f" ]; then
    echo "missing canonical file: $f" >&2
    exit 1
  fi
  sha=$(hash_file "$f")
  JSON=$(echo "$JSON" | jq --arg k "$f" --arg v "$sha" '. + {($k): $v}')
done

echo "$JSON" | jq -S .
