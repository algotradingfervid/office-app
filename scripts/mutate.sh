#!/usr/bin/env bash
# Mutation-test the code a story changed: scripts/mutate.sh <story-id> [base-ref]
# Uses gremlins (github.com/go-gremlins/gremlins v0.6.0) limited to lines changed since base (origin/main or main).
# Prints every LIVED mutant (a test ran the code but nothing failed) as file:line and exits 1 if there is one.
# NOT COVERED mutants (no test reaches the line, usually `if err != nil` returns) are listed for information.
# Raw JSON and log go to factory/evidence/<id>/raw/ (git-ignored). cmd/ and internal/testapp/ are skipped.
# Install once: go install github.com/go-gremlins/gremlins/cmd/gremlins@v0.6.0
set -euo pipefail
id="${1:?usage: scripts/mutate.sh <story-id> [base-ref]}"
cd "$(git rev-parse --show-toplevel)"
base="${2:-}"
if [ -z "$base" ]; then
  base="origin/main"
  git rev-parse --verify -q "$base" >/dev/null || base="main"
fi
command -v gremlins >/dev/null || { echo "gremlins not installed: go install github.com/go-gremlins/gremlins/cmd/gremlins@v0.6.0" >&2; exit 2; }
out="factory/evidence/$id/raw"
mkdir -p "$out"
changed="$(git diff --name-only "$base"...HEAD -- '*.go' | grep -v -e '_test\.go$' -e '^cmd/' -e '^internal/testapp/' || true)"
if [ -z "$changed" ]; then
  echo "mutation: no non-test Go files changed since $base"
  echo '{"files":[]}' >"$out/mutation.json"
  exit 0
fi
echo "mutation: base $base; changed files:"
echo "$changed" | sed 's/^/  /'
gremlins unleash --diff "$base" --exclude-files '^cmd/' --exclude-files '^internal/testapp/' \
  --output "$out/mutation.json" . >"$out/mutation.log" 2>&1 || true
python3 - "$out/mutation.json" <<'EOF'
import json, sys
data = json.load(open(sys.argv[1]))
rows = {"KILLED": [], "LIVED": [], "NOT COVERED": []}
for f in data.get("files", []):
    for m in f.get("mutations", []):
        if m.get("status") in rows:
            rows[m["status"]].append(f'{f["file_name"]}:{m["line"]}:{m["column"]}  {m["type"]}')
print(f'mutation: killed {len(rows["KILLED"])}, lived {len(rows["LIVED"])}, not covered {len(rows["NOT COVERED"])}')
for label in ("LIVED", "NOT COVERED"):
    for r in rows[label]:
        print(f"  {label:<11} {r}")
sys.exit(1 if rows["LIVED"] else 0)
EOF
