#!/usr/bin/env bash
# Push only go-codebuddy/ contents as the remote repository root.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

PREFIX="go-codebuddy"
BRANCH="publish/codebuddy"
REMOTE="${REMOTE:-origin}"
TARGET_REF="${TARGET_REF:-main}"
DRY_RUN=0
ALSO_MIRROR=0

for arg in "$@"; do
  case "$arg" in
    --dry-run) DRY_RUN=1 ;;
    --also-mirror) ALSO_MIRROR=1 ;;
    --remote=*) REMOTE="${arg#*=}" ;;
    --target=*) TARGET_REF="${arg#*=}" ;;
    -h|--help)
      cat <<'EOF'
Usage: push-codebuddy.sh [--dry-run] [--also-mirror] [--remote=origin] [--target=main]

Splits go-codebuddy/ into a temporary publish branch and force-pushes it
so the GitHub repo root is the CodeBuddy product (README, cmd, docs, …),
not this outer integrator workspace.
EOF
      exit 0
      ;;
    *)
      echo "unknown arg: $arg" >&2
      exit 2
      ;;
  esac
done

if [[ ! -d "$PREFIX" ]]; then
  echo "missing $PREFIX/" >&2
  exit 1
fi

echo "==> subtree split --prefix=$PREFIX -> $BRANCH"
git subtree split --prefix="$PREFIX" -b "$BRANCH"

echo "==> publish tree (top-level):"
git ls-tree --name-only "$BRANCH" | sed 's/^/  /'

if [[ "$DRY_RUN" -eq 1 ]]; then
  echo "==> dry-run: not pushing. Branch kept as $BRANCH"
  exit 0
fi

echo "==> push $BRANCH -> $REMOTE $TARGET_REF"
git push --force-with-lease "$REMOTE" "$BRANCH:$TARGET_REF"

if [[ "$ALSO_MIRROR" -eq 1 ]]; then
  echo "==> push $BRANCH -> mirror $TARGET_REF"
  git push --force-with-lease mirror "$BRANCH:$TARGET_REF"
fi

echo "==> done"
