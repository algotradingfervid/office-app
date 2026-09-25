#!/usr/bin/env python3
"""Read factory/stories/*.md and answer the line's questions.

  stories.py check [--all]   size, claim and collision rules for /slice (unmerged stories; --all audits every story)
  stories.py ready           counts, parallel waves, hotspots, STATE.md length
  stories.py stats           median hours per station and review rounds per lane
"""
import fnmatch, glob, json, os, re, statistics, subprocess, sys
from datetime import datetime

MAX_FILES, MAX_BODY, HOTSPOT, MAX_STATE = 5, 60, 3, 60
LANES, TIERS, KINDS = {"core", "full"}, {"core", "usable", "launch", "later"}, {"feature", "contract"}
ACTIVE = {"building", "proving", "review"}


def value(raw):
    raw = raw.strip()
    if raw.startswith("["):
        inner = raw[: raw.rindex("]") + 1]
        try:
            items = json.loads(inner)
        except json.JSONDecodeError:
            items = inner[1:-1].split(",")
        return [str(i).strip().strip("'\"") for i in items if str(i).strip().strip("'\"")]
    if raw.startswith('"'):
        return raw[1 : raw.index('"', 1)]
    return raw.split(" #")[0].strip()


def story_id(v):
    v = str(v).strip()
    return v.zfill(3) if v.isdigit() else v


def load():
    stories = []
    for path in sorted(glob.glob("factory/stories/*.md")):
        text = open(path, encoding="utf-8").read()
        m = re.match(r"---\n(.*?)\n---\n?(.*)", text, re.S)
        if not m:
            continue
        meta, key = {}, None
        for line in m.group(1).splitlines():
            item = re.match(r"\s+-\s+(.*)", line)
            if item and key and isinstance(meta.get(key), list):
                meta[key].append(item.group(1).split(" #")[0].strip().strip("'\""))
                continue
            k, sep, v = line.partition(":")
            if sep and re.fullmatch(r"[a-z-]+", k):
                key = k
                try:
                    meta[k] = value(v) if v.split(" #")[0].strip() else []
                except ValueError:
                    meta[k] = v.strip()
        for k in ("files", "needs"):
            if not isinstance(meta.get(k, []), list):
                meta[k] = [meta[k]] if meta[k] else []
        meta["needs"] = [story_id(n) for n in meta.get("needs", [])]
        meta["id"] = story_id(meta.get("id") or os.path.basename(path)[:3])
        meta["path"], meta["body_lines"] = path, len(m.group(2).strip().splitlines())
        meta.setdefault("files", []), meta.setdefault("needs", [])
        stories.append(meta)
    return stories


def is_test(f):
    return bool(re.search(r"(_test\.|\.test\.|\.spec\.|(^|/)test_|(^|/)tests?/|testdata/)", f)) or f.startswith("the tests")


def overlap(a, b):
    return any(x == y or fnmatch.fnmatch(x, y) or fnmatch.fnmatch(y, x) for x in a for y in b)


def ancestors(s, by_id, seen=None):
    seen = set() if seen is None else seen
    for n in s["needs"]:
        if n not in seen and n in by_id:
            seen.add(n)
            ancestors(by_id[n], by_id, seen)
    return seen


def check(all_stories):
    stories = load()
    by_id = {s["id"]: s for s in stories}
    todo = [s for s in stories if all_stories or s.get("status") != "merged"]
    fails = []
    for s in todo:
        sid, files = s["id"], s["files"]
        code = [f for f in files if not is_test(f)]
        if not s.get("check"):
            fails.append((sid, "no check"))
        if s.get("lane") not in LANES:
            fails.append((sid, f"lane must be one of {sorted(LANES)}"))
        if s.get("tier") not in TIERS:
            fails.append((sid, f"tier must be one of {sorted(TIERS)}"))
        if s.get("kind", "feature") not in KINDS:
            fails.append((sid, f"kind must be one of {sorted(KINDS)}"))
        if not files:
            fails.append((sid, "no file claim"))
        if len(code) > MAX_FILES:
            fails.append((sid, f"claims {len(code)} non-test files (max {MAX_FILES}): split it"))
        for f in files:
            if f in ("*", "**") or f.endswith("/*") or f.endswith("/**"):
                fails.append((sid, f"claims a whole directory: {f}"))
        if s["body_lines"] > MAX_BODY:
            fails.append((sid, f"body is {s['body_lines']} lines (max {MAX_BODY}): say what, not how"))
        for n in s["needs"]:
            if n not in by_id:
                fails.append((sid, f"needs unknown story {n}"))
    for i, a in enumerate(todo):
        for b in todo[i + 1 :]:
            related = b["id"] in ancestors(a, by_id) or a["id"] in ancestors(b, by_id)
            if not related and overlap(a["files"], b["files"]):
                fails.append((f"{a['id']}+{b['id']}", "claim the same files but neither needs the other"))
    for sid, msg in fails:
        print(f"FAIL {sid}: {msg}")
    print(f"{len(todo)} stories checked, {len(fails)} failures")
    return 1 if fails else 0


def ready():
    stories = load()
    for st in ["ready", "building", "proving", "review", "blocked", "merged", "backlog"]:
        print(f"{st:<9} {sum(1 for s in stories if s.get('status') == st)}")
    busy = [f for s in stories if s.get("status") in ACTIVE for f in s["files"]]
    pool = [s for s in stories if s.get("status") == "ready"]
    waves = []
    while pool:
        wave, taken, rest = [], list(busy), []
        for s in pool:
            (rest if overlap(s["files"], taken) else wave).append(s)
            if s in wave:
                taken += s["files"]
        if not wave:
            break
        waves.append(wave)
        pool = rest
    print("---")
    for n, wave in enumerate(waves, 1):
        print(f"wave {n}{' (start now)' if n == 1 else ' (after wave ' + str(n - 1) + ' merges)'}:")
        for s in wave:
            print(f"  /isolate {s['id']}   [{s.get('lane', '?')}] {s.get('title', '')}")
    if pool:
        print("waiting on files in flight: " + ", ".join(s["id"] for s in pool))
    try:
        cap = int(re.search(r"agents-recommended:\s*(\d+)", open("factory/STATE.md").read()).group(1))
    except (OSError, AttributeError):
        cap = 3
    free = max(0, cap - sum(1 for s in stories if s.get("status") in ACTIVE))
    print(f"agents: {min(free, len(waves[0]) if waves else 0)} can start now (cap {cap}, set in STATE.md)")
    counts = {}
    for s in stories:
        for f in s["files"]:
            if not is_test(f):
                counts[f] = counts.get(f, 0) + 1
    hot = sorted(((c, f) for f, c in counts.items() if c >= HOTSPOT), reverse=True)
    if hot:
        print("--- hotspots (claimed by >= %d stories; give the module its own registration or split the file):" % HOTSPOT)
        for c, f in hot[:10]:
            print(f"  {c:>3}  {f}")
    try:
        n = sum(1 for _ in open("factory/STATE.md"))
        if n > MAX_STATE:
            print(f"--- WARNING: factory/STATE.md is {n} lines (max {MAX_STATE}); move history to factory/HISTORY.md")
    except OSError:
        pass
    return 0


def when(v):
    for fmt in ("%Y-%m-%dT%H:%M", "%Y-%m-%d %H:%M", "%Y-%m-%d"):
        try:
            return datetime.strptime(str(v).strip()[:16], fmt)
        except ValueError:
            continue
    return None


def stats():
    stories = load()
    steps = [("started", "built"), ("built", "proved"), ("proved", "reviewed"), ("reviewed", "merged"), ("started", "merged")]
    for a, b in steps:
        hours = [(when(s[b]) - when(s[a])).total_seconds() / 3600 for s in stories
                 if when(s.get(a)) and when(s.get(b)) and when(s[b]) >= when(s[a])]
        med = f"{statistics.median(hours):.1f} h" if hours else "no data"
        print(f"{a:>8} -> {b:<8} median {med}  (n={len(hours)})")
    for lane in sorted(LANES):
        rounds = [int(s["review-rounds"]) for s in stories if s.get("lane") == lane and s.get("status") == "merged" and str(s.get("review-rounds", "")).isdigit()]
        if rounds:
            print(f"lane {lane:<5} review rounds: median {statistics.median(rounds)}, max {max(rounds)}  (n={len(rounds)})")
    return 0


if __name__ == "__main__":
    os.chdir(subprocess.run(["git", "rev-parse", "--show-toplevel"], capture_output=True, text=True).stdout.strip() or ".")
    cmd = sys.argv[1] if len(sys.argv) > 1 else "ready"
    sys.exit({"check": lambda: check("--all" in sys.argv), "ready": ready, "stats": stats}.get(cmd, lambda: print(__doc__) or 2)())
