#!/usr/bin/env bash
# Capture the evidence for a story: scripts/prove.sh <story-id>
# Run from the story's worktree. Writes factory/evidence/<id>/report.md (verdict left for /prove to set),
# after/ (and before/ for screen stories) screenshots, and raw logs under raw/ (git-ignored).
#   1. make check                     (must pass)
#   2. the story's ```check commands  (shell lines in the story's Check section; each must exit 0)
#   3. screen stories: the ```journey block via scripts/journey.py, on origin/main (before) and this branch (after)
#   4. scripts/mutate.sh              (no LIVED mutants)
set -uo pipefail
id="${1:?usage: scripts/prove.sh <story-id>}"
root="$(git rev-parse --show-toplevel)"
cd "$root"
story="$(ls factory/stories/"$id"-*.md 2>/dev/null | head -n1)"
[ -n "$story" ] || { echo "no story file for $id" >&2; exit 2; }
ev="factory/evidence/$id"
mkdir -p "$ev/raw"
base="origin/main"; git rev-parse --verify -q "$base" >/dev/null || base="main"
screen="$(sed -n 's/^screen:[[:space:]]*//p' "$story" | head -n1 | sed 's/[[:space:]]*#.*//')"
check_line="$(sed -n 's/^check:[[:space:]]*//p' "$story" | head -n1)"
code_commit="$(git log -1 --format=%H -- . ':!factory/evidence')"
status=0
section() { printf '\n## %s\n\n' "$1" >>"$ev/report.md"; }
block() { printf '```\n%s\n```\n' "$1" >>"$ev/report.md"; }

cat >"$ev/report.md" <<EOF
# Evidence — story $id

story: $story
code-commit: $code_commit
base: $base ($(git rev-parse --short "$base"))
captured: $(date '+%Y-%m-%d %H:%M')
check: $check_line
verdict: PENDING   # /prove sets PASS or GAP after reading this report
EOF

section "1. make check"
if make check >"$ev/raw/make-check.log" 2>&1; then r="PASS"; else r="FAIL"; status=1; fi
block "\$ make check   -> $r
$(grep -E '^(ok|FAIL|---|imports:|[0-9]+ issues)|BOUNDARY|gofmt needed' "$ev/raw/make-check.log" | tail -40)"

checks="$(awk '/^```check$/{f=1;next} /^```$/{f=0} f' "$story")"
if [ -n "$checks" ]; then
  section "2. Story check commands"
  while IFS= read -r cmd; do
    [ -z "$cmd" ] && continue
    out="$(bash -c "$cmd" 2>&1)"; rc=$?
    [ $rc -ne 0 ] && status=1
    block "\$ $cmd   -> exit $rc
$(printf '%s\n' "$out" | tail -30)"
  done <<<"$checks"
fi

if [ -n "$screen" ] && [ "$screen" != "none" ] && ! grep -q '^none' <<<"$screen"; then
  section "3. Journey (390 px and 1280 px, light)"
  mkdir -p .preview
  CGO_ENABLED=0 go build -o ".preview/prove-$id-after" ./cmd/officeapp
  if python3 scripts/journey.py "$story" --bin ".preview/prove-$id-after" --out "$ev/after" >"$ev/raw/journey-after.log" 2>&1; then r="PASS"; else r="FAIL"; status=1; fi
  block "after ($(git rev-parse --short HEAD)) -> $r
$(cat "$ev/raw/journey-after.log")"
  bt=".preview/prove-$id-before-src"
  git worktree remove --force "$bt" >/dev/null 2>&1 || true
  if git worktree add --detach "$bt" "$base" >/dev/null 2>&1; then
    (cd "$bt" && CGO_ENABLED=0 go build -o "$root/.preview/prove-$id-before" ./cmd/officeapp) >/dev/null 2>&1
    python3 scripts/journey.py "$story" --bin ".preview/prove-$id-before" --out "$ev/before" >"$ev/raw/journey-before.log" 2>&1
    block "before ($(git rev-parse --short "$base")), expected to fail where the feature is new:
$(cat "$ev/raw/journey-before.log")"
    git worktree remove --force "$bt" >/dev/null 2>&1
  fi
  printf '\nScreenshots: `after/` (this branch), `before/` (%s).\n' "$base" >>"$ev/report.md"
fi

section "4. Mutation (scripts/mutate.sh $id)"
if scripts/mutate.sh "$id" "$base" >"$ev/raw/mutate.out" 2>&1; then r="PASS"; else r="FAIL"; status=1; fi
block "\$ scripts/mutate.sh $id   -> $r
$(cat "$ev/raw/mutate.out")"

echo "evidence: $ev/report.md  (automatic checks: $([ $status = 0 ] && echo PASS || echo FAIL))"
exit $status
