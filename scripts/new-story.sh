#!/usr/bin/env bash
# Create an isolated worktree for a story: scripts/new-story.sh <story-id>
set -euo pipefail
id="${1:?usage: scripts/new-story.sh <story-id>}"
story="$(ls factory/stories/"${id}"-*.md 2>/dev/null | head -n1 || true)"
[ -n "$story" ] || { echo "no story file for id $id in factory/stories/" >&2; exit 1; }
slug="$(basename "$story" .md | sed "s/^${id}-//")"
branch="story/${id}-${slug}"
dir=".claude/worktrees/story-${id}"

git fetch origin main --quiet 2>/dev/null || true
base="origin/main"; git rev-parse --verify -q "$base" >/dev/null || base="main"

if [ -d "$dir" ]; then echo "worktree exists: $dir"; exit 0; fi
mkdir -p .claude/worktrees
git worktree add -b "$branch" "$dir" "$base" >/dev/null

# Copy untracked files the app needs (one path per line in .worktreeinclude)
if [ -f .worktreeinclude ]; then
  while IFS= read -r f; do
    [ -z "$f" ] && continue
    [ -e "$f" ] && mkdir -p "$dir/$(dirname "$f")" && cp -r "$f" "$dir/$f"
  done < .worktreeinclude
fi

# Install deps if a known manifest exists
( cd "$dir"
  if   [ -f pnpm-lock.yaml ]; then pnpm install --frozen-lockfile --silent
  elif [ -f package-lock.json ]; then npm ci --silent
  elif [ -f yarn.lock ]; then yarn install --silent
  elif [ -f pyproject.toml ] && command -v uv >/dev/null; then uv sync --quiet
  elif [ -f requirements.txt ]; then python -m venv .venv && .venv/bin/pip install -q -r requirements.txt
  fi
) || echo "dependency install failed; continue manually" >&2

echo "$dir"
echo "branch: $branch (from $base)"
echo "open a fresh session there:  claude -w story-${id}"
