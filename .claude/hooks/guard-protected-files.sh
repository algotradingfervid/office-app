#!/usr/bin/env bash
# PreToolUse guard for Edit/Write: block edits to protected paths.
# Exit 2 = block (stderr is shown to the agent). Exit 0 = allow.
# /tool-up and /learn append product-specific patterns to PROTECTED below.
set -euo pipefail
input="$(cat)"
path="$(printf '%s' "$input" | sed -n 's/.*"file_path"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)"
[ -z "$path" ] && exit 0

PROTECTED=(
  ".env" ".env.*" "*.pem" "*.key" "secrets/*"
  "factory/SPEC.md"            # changes only via /define
  "factory/BRAND.md"           # changes only via /brand
  "factory/design/tokens/*"    # generated from tokens.json
  "factory/evidence/*"         # never edit evidence; recapture
  "*.lock" "package-lock.json" "pnpm-lock.yaml" "yarn.lock" "Cargo.lock" "poetry.lock"
  # product (added by /tool-up)
  "go.sum"                                 # regenerate with `go mod tidy`
  "internal/core/web/static/htmx.min.js"   # vendored, pinned; see static/VERSIONS
  "internal/core/web/static/pico.min.css"  # vendored, pinned; see static/VERSIONS
)

# Product (added by /tool-up): a migration that is already on main is history; never edit it, add a new one.
case "$(basename "$path")" in
  [0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]_*.go)
    rel="${path#"$(git rev-parse --show-toplevel 2>/dev/null)/"}"
    if git cat-file -e "main:$rel" 2>/dev/null; then
      echo "SHIPSHOW guard: '$rel' is a merged migration. Add a new timestamp-named migration instead." >&2; exit 2
    fi ;;
esac

# Allow-list: the station that owns a file writes its name to .claude/station while it runs.
station="$(cat .claude/station 2>/dev/null || true)"
case "$path" in
  */factory/SPEC.md|factory/SPEC.md)   [ "$station" = "define" ] && exit 0 ;;
  */factory/BRAND.md|factory/BRAND.md) [ "$station" = "brand" ]  && exit 0 ;;
esac

for pat in "${PROTECTED[@]}"; do
  # shellcheck disable=SC2254
  case "$path" in
    $pat|*/$pat) echo "SHIPSHOW guard: '$path' is protected ($pat). Use the owning station, or ask the human." >&2; exit 2 ;;
  esac
done
exit 0
