#!/usr/bin/env bash
# After a story's PR is merged: remove worktree + branch, mark merged, promote newly-unblocked stories to ready.
set -euo pipefail
id="${1:?usage: scripts/finish-story.sh <story-id>}"
story="$(ls factory/stories/"${id}"-*.md | head -n1)"
dir=".claude/worktrees/story-${id}"
branch="$(git -C "$dir" branch --show-current 2>/dev/null || git branch --list "story/${id}-*" | tr -d ' *' | head -n1)"

git worktree remove --force "$dir" 2>/dev/null || true
[ -n "$branch" ] && git branch -D "$branch" 2>/dev/null || true
git fetch origin main --quiet 2>/dev/null || true

now="$(date '+%Y-%m-%d %H:%M')"
sed -i.bak -E -e 's/^status:.*/status: merged/' -e "s/^merged:.*/merged: $now/" "$story" && rm -f "$story.bak"
merged=$(grep -l '^status: merged' factory/stories/*.md | sed -E 's#.*/([0-9]+)-.*#\1#' | sort -u)

promoted=()
for f in factory/stories/*.md; do
  grep -q '^status: backlog' "$f" || continue
  needs=$(sed -n 's/^needs:[[:space:]]*\[\(.*\)\].*/\1/p' "$f" | tr -d '" ' | tr ',' '\n' | sed '/^$/d')
  ok=1
  for n in $needs; do echo "$merged" | grep -qx "$n" || { ok=0; break; }; done
  if [ $ok -eq 1 ]; then sed -i.bak -E 's/^status:.*/status: ready/' "$f" && rm -f "$f.bak"; promoted+=("$(basename "$f" .md)"); fi
done
echo "story $id merged."
[ ${#promoted[@]} -gt 0 ] && echo "now ready: ${promoted[*]}" || echo "nothing newly unblocked."
