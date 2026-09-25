#!/usr/bin/env bash
# PreToolUse guard for Bash: block commits/pushes on main unless the office is open.
# Office stations write their name to .claude/station while they run (see CLAUDE.md).
set -euo pipefail
input="$(cat)"
cmd="$(printf '%s' "$input" | sed -n 's/.*"command"[[:space:]]*:[[:space:]]*"\(.*\)".*/\1/p' | head -n1)"
[ -z "$cmd" ] && exit 0

case "$cmd" in
  *"git commit"*|*"git push"*|*"git merge"*)
    branch="$(git branch --show-current 2>/dev/null || true)"
    if [ "$branch" = "main" ] || [ "$branch" = "master" ]; then
      case "$(cat .claude/station 2>/dev/null || true)" in
        blueprint|tool-up|shipshow-update|slice|isolate|learn|spike) exit 0 ;;
        *) echo "SHIPSHOW guard: you are on '$branch'. Stories live in worktrees (/isolate). Only office stations commit to main." >&2; exit 2 ;;
      esac
    fi
    ;;
esac
case "$cmd" in
  *"--force "*|*" -f "*"push"*|*"push -f"*|*"push --force"*)
    echo "SHIPSHOW guard: force push is blocked. Use --force-with-lease on a story branch, never on main." >&2; exit 2 ;;
esac
exit 0
