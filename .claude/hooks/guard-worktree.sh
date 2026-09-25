#!/usr/bin/env bash
# PreToolUse guard for Edit/Write/MultiEdit: inside a story worktree, block writes to any path outside it.
# A story agent that writes into main's checkout by absolute path leaves main broken (lesson, story 023).
# Temp folders and ~/.claude (scratch, plans, memory) stay writable.
set -euo pipefail
python3 -c '
import json, os, subprocess, sys
d = json.load(sys.stdin)
cwd = d.get("cwd") or os.getcwd()
path = (d.get("tool_input") or {}).get("file_path") or ""
if not path:
    sys.exit(0)
top = subprocess.run(["git", "-C", cwd, "rev-parse", "--show-toplevel"], capture_output=True, text=True).stdout.strip()
if "/.claude/worktrees/story-" not in top + "/":
    sys.exit(0)
target = os.path.realpath(os.path.join(cwd, path))
root = os.path.realpath(top)
common = subprocess.run(["git", "-C", cwd, "rev-parse", "--path-format=absolute", "--git-common-dir"], capture_output=True, text=True).stdout.strip()
main = os.path.dirname(os.path.realpath(common)) if common else root
under = lambda a: target == a or target.startswith(a.rstrip(os.sep) + os.sep)
scratch = [os.path.realpath(p) for p in ("/tmp", os.environ.get("TMPDIR", "/tmp"), os.path.expanduser("~/.claude"))]
if not under(root) and (under(main) or not any(under(a) for a in scratch)):
    print(f"SHIPSHOW guard: this session builds in {root}. {target} is outside the story worktree; write inside it.", file=sys.stderr)
    sys.exit(2)
'
