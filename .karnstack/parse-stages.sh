#!/usr/bin/env bash
# Parse `go test -v` output and emit JSON:
#
#   {"passing_stages": [1, 2, 3]}
#
# A stage passes only if every top-level TestStageNN_* function for that
# stage reports PASS. Any FAIL on a TestStageNN_* function fails that
# stage. Subtest results (lines containing a slash like
# `TestStage01_Foo/case_a`) are ignored - the table-driven subtest
# pattern is allowed, and the top-level result already aggregates them.
#
# Portable across bash 3.2 (macOS) and bash 4+ (Linux). No associative
# arrays; uses sort + comm to compute set difference.
#
# Usage:
#   bash parse-stages.sh test-output.txt
#   mise run all | bash parse-stages.sh /dev/stdin

set -euo pipefail

LOG="${1:-/dev/stdin}"
TMP=$(mktemp -t parse-stages.XXXXXX)
trap 'rm -f "$TMP"' EXIT

# Extract one line per top-level test result: "<RESULT> <NN>".
# The trailing literal space before `(` ensures we match a top-level
# function name with no slash, since subtests embed a `/` in the name.
grep -E '^--- (PASS|FAIL): TestStage[0-9]{2}_[A-Za-z0-9_]+ \(' "$LOG" \
  | sed -E 's/^--- (PASS|FAIL): TestStage([0-9]{2})_.* \(.*/\1 \2/' \
  > "$TMP" || true

FAILED=$(awk '$1=="FAIL" {print $2}' "$TMP" | sort -u)
SEEN=$(awk '{print $2}' "$TMP" | sort -u)

if [ -z "$SEEN" ]; then
  echo '{"passing_stages": []}'
  exit 0
fi

PASSING=$(comm -23 <(printf '%s\n' "$SEEN") <(printf '%s\n' "$FAILED"))

if [ -z "$PASSING" ]; then
  echo '{"passing_stages": []}'
  exit 0
fi

# Strip leading zeros and emit a sorted JSON array of integers.
printf '%s\n' "$PASSING" \
  | sed -E 's/^0+//' \
  | sort -n \
  | jq -R 'tonumber' \
  | jq -sc '{passing_stages: .}'
