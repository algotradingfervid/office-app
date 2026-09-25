#!/usr/bin/env bash
# PostToolUse: format the file that was just edited, if a formatter exists for it.
# /tool-up extends the case list for the product's languages.
set -uo pipefail
input="$(cat)"
path="$(printf '%s' "$input" | sed -n 's/.*"file_path"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)"
[ -z "$path" ] || [ ! -f "$path" ] && exit 0
case "$path" in
  *.js|*.jsx|*.ts|*.tsx|*.css|*.json|*.md|*.html)
    command -v prettier >/dev/null 2>&1 && prettier --write --log-level silent "$path" >/dev/null 2>&1 ;;
  *.py)
    command -v ruff >/dev/null 2>&1 && ruff format -q "$path" >/dev/null 2>&1 ;;
  *.go)
    command -v gofmt >/dev/null 2>&1 && gofmt -w "$path" ;;
  *.rs)
    command -v rustfmt >/dev/null 2>&1 && rustfmt "$path" ;;
  *.swift)
    command -v swiftformat >/dev/null 2>&1 && swiftformat "$path" >/dev/null 2>&1 ;;
  *.kt)
    command -v ktlint >/dev/null 2>&1 && ktlint -F "$path" >/dev/null 2>&1 ;;
esac
exit 0
