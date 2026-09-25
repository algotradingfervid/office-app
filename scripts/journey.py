#!/usr/bin/env python3
"""The one journey runner: drive a built officeapp through the steps in a story's Check section.

Usage:
  scripts/journey.py STORY_FILE --bin PATH --out DIR [--widths 390,1280]

For each width it starts the binary on a fresh seeded database, runs the story's ```journey block in
Chromium (light theme) and stops the server. Screenshots land in DIR/<name>-<width>.png, and
DIR/journey.json records every step. Exit 1 if a step fails.

Journey steps, one per line (blank lines and # comments ignored):
  login E001                 sign in as a demo user (password from internal/core/auth/seed.go)
  goto /path                 open a page
  fill SELECTOR VALUE        type into a field (VALUE is the rest of the line)
  select SELECTOR VALUE      choose an <option> by value or label
  click SELECTOR             click (Playwright selectors: css, text=Submit, role=button[name="Save"])
  see TEXT                   the page shows TEXT
  notsee TEXT                the page does not show TEXT
  wait MS                    pause (use sparingly; htmx swaps are awaited automatically)
  shot NAME                  full-page screenshot
"""
import argparse, json, os, re, socket, subprocess, sys, tempfile, time, urllib.request
from pathlib import Path

try:
    from playwright.sync_api import sync_playwright
except ImportError:
    sys.exit("journey.py needs Playwright: pip install playwright && playwright install chromium")

DEMO_PASSWORD = "demo-pass-2026"


def steps_from(story):
    text = Path(story).read_text(encoding="utf-8")
    m = re.search(r"```journey\n(.*?)```", text, re.S)
    if not m:
        sys.exit(f"{story}: no ```journey block in the Check section")
    out = []
    for line in m.group(1).splitlines():
        line = line.strip()
        if line and not line.startswith("#"):
            verb, _, rest = line.partition(" ")
            out.append((verb, rest.strip()))
    return out


def free_port():
    s = socket.socket()
    s.bind(("127.0.0.1", 0))
    port = s.getsockname()[1]
    s.close()
    return port


def start(binary, data):
    subprocess.run([binary, "seed", "--dir", data], check=True, capture_output=True)
    port = free_port()
    log = open(os.path.join(data, "server.log"), "w")
    proc = subprocess.Popen([binary, "serve", "--dir", data, "--http", f"127.0.0.1:{port}"], stdout=log, stderr=log)
    url = f"http://127.0.0.1:{port}"
    for _ in range(100):
        try:
            urllib.request.urlopen(url + "/api/health", timeout=1)
            return proc, url
        except Exception:
            time.sleep(0.1)
    proc.kill()
    sys.exit(f"server did not start; see {data}/server.log")


def run(page, url, steps, out, width, record):
    for verb, arg in steps:
        entry = {"width": width, "step": f"{verb} {arg}".strip(), "ok": True}
        try:
            if verb == "login":
                page.goto(url + "/login")
                page.fill("input[name=identity]", arg)
                page.fill("input[name=password]", DEMO_PASSWORD)
                page.click("button[type=submit]")
                page.wait_for_load_state("networkidle")
            elif verb == "goto":
                page.goto(url + arg)
            elif verb in ("fill", "select"):
                sel, _, val = arg.partition(" ")
                (page.fill if verb == "fill" else page.select_option)(sel, val)
            elif verb == "click":
                page.click(arg)
            elif verb == "see":
                page.get_by_text(arg).first.wait_for(state="visible", timeout=5000)
            elif verb == "notsee":
                if page.get_by_text(arg).count() and page.get_by_text(arg).first.is_visible():
                    raise AssertionError(f"'{arg}' is visible")
            elif verb == "wait":
                page.wait_for_timeout(int(arg))
            elif verb == "shot":
                path = out / f"{arg}-{width}.png"
                page.screenshot(path=str(path), full_page=True)
                entry["file"] = path.name
            else:
                raise ValueError(f"unknown step '{verb}'")
            page.wait_for_load_state("networkidle")
        except Exception as e:  # record and stop this width
            entry.update(ok=False, error=str(e).splitlines()[0])
            record.append(entry)
            return False
        record.append(entry)
    return True


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("story")
    ap.add_argument("--bin", required=True)
    ap.add_argument("--out", required=True)
    ap.add_argument("--widths", default="390,1280")
    a = ap.parse_args()
    steps, out = steps_from(a.story), Path(a.out)
    out.mkdir(parents=True, exist_ok=True)
    record, ok = [], True
    with sync_playwright() as p:
        browser = p.chromium.launch()
        for width in [int(w) for w in a.widths.split(",")]:
            with tempfile.TemporaryDirectory() as data:
                proc, url = start(os.path.abspath(a.bin), data)
                try:
                    page = browser.new_page(viewport={"width": width, "height": 800}, color_scheme="light")
                    errors = []
                    page.on("console", lambda m: errors.append(m.text) if m.type == "error" else None)
                    ok = run(page, url, steps, out, width, record) and ok
                    if errors:
                        record.append({"width": width, "step": "console", "ok": False, "error": "; ".join(errors)})
                        ok = False
                    page.close()
                finally:
                    proc.terminate()
                    proc.wait()
        browser.close()
    (out / "journey.json").write_text(json.dumps(record, indent=2))
    for r in record:
        print(("PASS " if r["ok"] else "FAIL ") + f'[{r["width"]}] {r["step"]}' + ("" if r["ok"] else f'  -> {r["error"]}'))
    sys.exit(0 if ok else 1)


if __name__ == "__main__":
    main()
