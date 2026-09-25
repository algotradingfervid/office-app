#!/usr/bin/env bash
# Stories by status, parallel waves that can start now, file hotspots, STATE.md length.
#   scripts/ready.sh          the floor
#   scripts/ready.sh --stats  median hours per station and review rounds per lane
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
if [ "${1:-}" = "--stats" ]; then exec python3 scripts/stories.py stats; fi
exec python3 scripts/stories.py ready
