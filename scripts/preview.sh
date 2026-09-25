#!/usr/bin/env bash
# Build the current checkout and serve it on a free local port with fresh demo data.
#   scripts/preview.sh          start (prints the URL); a previous preview of this checkout is stopped first
#   scripts/preview.sh --stop   stop this checkout's preview
# Demo users: E001 (employee), E002 and E003 (approvers), E004 (HR admin); password in internal/core/auth/seed.go.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
name="$(git branch --show-current 2>/dev/null || true)"
name="${name:-detached}"
name="${name//\//-}"
dir=".preview/$name"
if [ -f "$dir/pid" ]; then
  kill "$(cat "$dir/pid")" 2>/dev/null || true
  rm -f "$dir/pid"
fi
if [ "${1:-}" = "--stop" ]; then echo "stopped $name"; exit 0; fi
rm -rf "$dir"
mkdir -p "$dir"
CGO_ENABLED=0 go build -o "$dir/officeapp" ./cmd/officeapp
"$dir/officeapp" seed --dir "$dir/pb_data" >/dev/null
port="$(python3 -c 'import socket; s=socket.socket(); s.bind(("127.0.0.1",0)); print(s.getsockname()[1])')"
nohup "$dir/officeapp" serve --dir "$dir/pb_data" --http "127.0.0.1:$port" >"$dir/server.log" 2>&1 &
echo $! >"$dir/pid"
for _ in $(seq 1 50); do
  if curl -fs "http://127.0.0.1:$port/api/health" >/dev/null 2>&1; then
    echo "http://127.0.0.1:$port"
    exit 0
  fi
  sleep 0.2
done
echo "preview failed to start; see $dir/server.log" >&2
exit 1
