#!/usr/bin/env bash
# Copy only what a browser renders (html, css, js, images, fonts) into a fresh temp folder outside the repo,
# keeping relative paths, so reviewer-blind sees the output and none of the documents behind it.
# Usage: scripts/blind-copy.sh <root> <path under root>...   e.g. scripts/blind-copy.sh factory/design wireframes components.css tokens icons.svg
# Prints the temp folder.
set -euo pipefail
root="${1:?usage: scripts/blind-copy.sh <root> <path under root>...}"
shift
[ "$#" -gt 0 ] || { echo "blind-copy: name at least one path under $root" >&2; exit 1; }
for p in "$@"; do
  [ -e "$root/$p" ] || { echo "blind-copy: $root/$p does not exist" >&2; exit 1; }
done
tmp="${TMPDIR:-/tmp}"
dest="$(mktemp -d "${tmp%/}/blind-review.XXXXXX")"
for p in "$@"; do
  if [ -d "$root/$p" ]; then
    mkdir -p "$dest/$p"
    rsync -a --prune-empty-dirs --include='*/' \
      --include='*.html' --include='*.css' --include='*.js' --include='*.svg' --include='*.png' --include='*.jpg' \
      --include='*.jpeg' --include='*.webp' --include='*.gif' --include='*.woff' --include='*.woff2' --include='*.ttf' \
      --exclude='*' "$root/$p/" "$dest/$p/"
  elif [ -f "$root/$p" ]; then
    mkdir -p "$dest/$(dirname "$p")"
    cp "$root/$p" "$dest/$p"
  fi
done
echo "$dest"
