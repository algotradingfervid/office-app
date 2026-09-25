#!/usr/bin/env bash
# PreToolUse guard for Bash (product, added by /learn 2026-09-25): block writing Go or template files inside a
# story worktree through a shell heredoc or a script. Edit/Write run the formatter and the other guards;
# shell writes skip them. Seen three times (tool-up trial, 003, 010), so the rule became this hook.
set -euo pipefail
python3 -c '
import json, re, sys
cmd = (json.load(sys.stdin).get("tool_input") or {}).get("command", "")
if ".claude/worktrees/story-" not in cmd and "/worktrees/story-" not in cmd:
    sys.exit(0)
writes_code = re.search(r"\.(go|html)\b", cmd) and (
    "<<" in cmd or re.search(r"open\([^)]*[\"\x27]w", cmd) or re.search(r">\s*\S+\.(go|html)\b", cmd)
    or re.search(r"\bsed\s+-i", cmd) or "write_text" in cmd)
if writes_code:
    print("SHIPSHOW guard: write Go/template files in a story worktree with Write/Edit, not a shell script or heredoc "
          "(hooks and gofmt only see Write/Edit). Reading, running tests and git are fine.", file=sys.stderr)
    sys.exit(2)
'
