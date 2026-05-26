#!/usr/bin/env bash
# .karnstack/sync.sh - pull canonical files from the karnstack template.
#
# Run this whenever karnstack ships an update to the workflow, helper
# scripts, mise config, or test files and your push starts failing with
# `files_modified` from /api/v1/byox/verify.
#
# Safe by design: only canonical files are touched. Your implementation
# (bloom/bloom.go) and any non-canonical edits are left alone. The script
# aborts if you have local uncommitted changes to any canonical file, so
# nothing of yours gets clobbered without you noticing.
#
# Usage:
#   bash .karnstack/sync.sh
#   mise run sync                   # same, via mise
#
# Override the upstream for testing:
#   TEMPLATE_OVERRIDE=acme/byox-bloom-filter-go bash .karnstack/sync.sh

set -euo pipefail

TEMPLATE="${TEMPLATE_OVERRIDE:-karnstack/byox-bloom-filter-go}"
REMOTE="karnstack-template"
BRANCH="main"

echo "syncing canonical files from ${TEMPLATE}@${BRANCH}"

if ! git remote get-url "$REMOTE" >/dev/null 2>&1; then
  git remote add "$REMOTE" "https://github.com/${TEMPLATE}.git"
fi
git fetch --quiet "$REMOTE" "$BRANCH"

# Files always considered canonical. The stage test files come from a
# remote glob below so newly-added test files come along automatically.
declare -a FILES=(
  ".github/workflows/verify-stages.yml"
  ".karnstack/compute-hashes.sh"
  ".karnstack/parse-stages.sh"
  ".karnstack/sync.sh"
  ".mise.toml"
)

# Discover stage test files from the upstream tree. Keeps sync robust to
# new stages landing in the template without code changes here. Uses a
# while-read loop (not mapfile) so this works on macOS's default
# bash 3.2.
REMOTE_TESTS=()
while IFS= read -r line; do
  [ -n "$line" ] && REMOTE_TESTS+=("$line")
done < <(
  git ls-tree -r --name-only "${REMOTE}/${BRANCH}" \
    | grep -E '^bloom/stage[0-9]+_[a-z0-9_]+_test\.go$' \
    || true
)
if [ "${#REMOTE_TESTS[@]}" -gt 0 ]; then
  FILES+=("${REMOTE_TESTS[@]}")
fi

# Refuse to overwrite uncommitted edits to a canonical file. The user
# should stash or commit deliberately before letting sync reset them.
DIRTY=$(git status --porcelain -- "${FILES[@]}" 2>/dev/null | awk '{print $2}')
if [ -n "$DIRTY" ]; then
  echo
  echo "ERROR: you have uncommitted changes to canonical files:"
  echo "$DIRTY" | sed 's/^/  /'
  echo
  echo "stash or commit them first, then re-run sync."
  exit 1
fi

for f in "${FILES[@]}"; do
  git checkout "${REMOTE}/${BRANCH}" -- "$f"
  echo "  synced $f"
done

# Surface (but don't delete) local stage tests that no longer exist
# upstream. A breaking template change might want them removed; we leave
# that to the user.
shopt -s nullglob
LOCAL_TESTS=( bloom/stage*_test.go )
shopt -u nullglob
for f in "${LOCAL_TESTS[@]}"; do
  found=0
  for r in "${REMOTE_TESTS[@]}"; do
    if [ "$r" = "$f" ]; then found=1; break; fi
  done
  if [ "$found" = "0" ]; then
    echo "  note: $f exists locally but not upstream; review whether to remove"
  fi
done

echo
echo "done. review the diff:"
echo "  git diff"
echo "then commit and push:"
echo "  git commit -am 'sync canonical files from karnstack'"
echo "  git push"
