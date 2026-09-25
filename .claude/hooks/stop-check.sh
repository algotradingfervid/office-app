#!/usr/bin/env bash
# Stop hook: on a story branch, refuse to end the session while `make check` fails.
# Exit 2 blocks the stop and feeds stderr back to the agent. Gives up after 8 consecutive blocks (Claude Code override).
set -uo pipefail
branch="$(git branch --show-current 2>/dev/null || true)"
case "$branch" in
  story/*) ;;
  *) exit 0 ;;
esac
[ -f Makefile ] && grep -q '^check:' Makefile || exit 0
if ! make check >/tmp/shipshow-check.log 2>&1; then
  echo "SHIPSHOW: make check is failing on $branch. Fix the root cause before stopping. Last lines:" >&2
  tail -n 30 /tmp/shipshow-check.log >&2
  exit 2
fi
exit 0
