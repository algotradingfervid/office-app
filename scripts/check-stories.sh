#!/usr/bin/env bash
# /slice gate: every unmerged story is small, claims exact files, and no two stories that can run together claim the same file.
#   scripts/check-stories.sh        unmerged stories
#   scripts/check-stories.sh --all  audit every story, merged included
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
exec python3 scripts/stories.py check "$@"
