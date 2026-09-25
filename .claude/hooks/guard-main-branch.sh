#!/usr/bin/env bash
# PreToolUse guard for Bash: block commits/pushes on main unless the office is open.
# Office stations write their name to .claude/station while they run (see CLAUDE.md).
set -euo pipefail
input="$(cat)"
cmd="$(printf '%s' "$input" | python3 -c 'import json,sys; print((json.load(sys.stdin).get("tool_input") or {}).get("command",""))' 2>/dev/null || true)"
[ -z "$cmd" ] && exit 0

# git commit / push / merge (not merge-base), including `git -C <dir> ...` forms
if printf '%s' "$cmd" | grep -Eq '(^|[;&|[:space:]])git([[:space:]]+-C[[:space:]]+[^[:space:]]+)?[[:space:]]+(commit|push|merge)([[:space:]]|$)'; then
  # The command may run elsewhere than the session: `cd <dir> && ...` or `git -C <dir> ...`.
  dir="$(printf '%s' "$cmd" | sed -n 's/^[[:space:]]*cd[[:space:]][[:space:]]*\([^;&|]*[^;&| ]\)[[:space:]]*&&.*/\1/p' | head -n1)"
  [ -z "$dir" ] && dir="$(printf '%s' "$cmd" | sed -n 's/.*git[[:space:]][[:space:]]*-C[[:space:]][[:space:]]*\([^[:space:]]*\).*/\1/p' | head -n1)"
  dir="$(printf '%s' "$dir" | tr -d "\"'")"
  branch="$(git -C "${dir:-.}" branch --show-current 2>/dev/null || git branch --show-current 2>/dev/null || true)"
  if [ "$branch" = "main" ] || [ "$branch" = "master" ]; then
    case "$(cat .claude/station 2>/dev/null || true)" in
      blueprint|tool-up|shipshow-update|slice|isolate|learn|spike) ;;
      *) echo "SHIPSHOW guard: you are on '$branch'. Stories live in worktrees (/isolate). Only office stations commit to main." >&2; exit 2 ;;
    esac
  fi
fi
case "$cmd" in
  *"--force "*|*" -f "*"push"*|*"push -f"*|*"push --force"*)
    echo "SHIPSHOW guard: force push is blocked. Use --force-with-lease on a story branch, never on main." >&2; exit 2 ;;
esac
exit 0
