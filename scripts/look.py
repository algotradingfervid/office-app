#!/usr/bin/env python3
"""Look at HTML pages the way a person would: in a real browser, at phone and desktop width, light and dark.

Usage:
  scripts/look.py PAGE_OR_DIR ... --out DIR [--root DIR] [--widths 390,1280] [--themes light,dark] [--crawl ENTRY]

Serves --root (default: the current folder) over local HTTP, so relative CSS, fonts and SVG sprites load as
they will in the app. For every page x width x theme it saves a full-page PNG in --out and writes
--out/look.json. With --crawl it follows local links from ENTRY, screenshots every page it reaches, and
reports links to missing files and given pages the entry cannot reach.

Exits 1 on console errors, assets that fail to load, horizontal scroll, missing link targets or unreachable
pages. Exit 0 does not mean it looks right: open the screenshots.

Themes are applied both ways a page may listen: prefers-color-scheme and <html data-theme="...">.
Needs: pip install playwright && playwright install chromium
"""
import argparse
import functools
import json
import sys
import threading
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from urllib.parse import unquote, urlparse

try:
    from playwright.sync_api import sync_playwright
except ImportError:
    sys.exit("look.py needs Playwright: pip install playwright && playwright install chromium")


class QuietHandler(SimpleHTTPRequestHandler):
    def log_message(self, *args):
        pass


def serve(root):
    handler = functools.partial(QuietHandler, directory=str(root))
    server = ThreadingHTTPServer(("127.0.0.1", 0), handler)
    threading.Thread(target=server.serve_forever, daemon=True).start()
    return server, f"http://127.0.0.1:{server.server_address[1]}"


def to_url(base, root, path):
    return f"{base}/{path.relative_to(root).as_posix()}"


def to_path(root, base, url):
    if not url.startswith(base + "/"):
        return None
    return root / unquote(urlparse(url).path).lstrip("/")


def crawl(page, root, base, entry):
    """Follow local links from entry. Returns (reached pages, [(from page, missing target)])."""
    reached, missing, queue = set(), [], [entry]
    while queue:
        current = queue.pop()
        if current in reached:
            continue
        reached.add(current)
        page.goto(to_url(base, root, current), wait_until="load")
        for href in page.eval_on_selector_all("a[href]", "els => els.map(e => e.href)"):
            target = to_path(root, base, href)
            if target is None:
                continue
            if not target.exists():
                missing.append((str(current.relative_to(root)), str(target.relative_to(root))))
            elif target.suffix == ".html" and target not in reached:
                queue.append(target)
    return reached, missing


def main():
    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("pages", nargs="*", help="HTML files, or folders (their *.html)")
    ap.add_argument("--out", required=True)
    ap.add_argument("--root", default=".")
    ap.add_argument("--widths", default="390,1280")
    ap.add_argument("--themes", default="light,dark")
    ap.add_argument("--crawl", metavar="ENTRY")
    args = ap.parse_args()

    root = Path(args.root).resolve()
    out = Path(args.out)
    out.mkdir(parents=True, exist_ok=True)
    pages = set()
    for p in map(Path, args.pages):
        pages.update(p.resolve() for p in (sorted(p.glob("*.html")) if p.is_dir() else [p]))
    if not pages and not args.crawl:
        sys.exit("look.py: give pages, a folder, or --crawl ENTRY")

    server, base = serve(root)
    report = {"shots": [], "missing_links": [], "unreachable": []}
    problems = 0
    with sync_playwright() as pw:
        browser = pw.chromium.launch()
        if args.crawl:
            page = browser.new_page()
            reached, missing = crawl(page, root, base, Path(args.crawl).resolve())
            page.close()
            report["missing_links"] = [{"from": f, "to": t} for f, t in missing]
            report["unreachable"] = sorted(str(p.relative_to(root)) for p in pages - reached)
            problems += len(missing) + len(report["unreachable"])
            pages |= reached

        for width in map(int, args.widths.split(",")):
            for theme in args.themes.split(","):
                ctx = browser.new_context(viewport={"width": width, "height": 844 if width < 768 else 900},
                                          color_scheme=theme)
                for path in sorted(pages):
                    page = ctx.new_page()
                    errors, failed = [], {}
                    page.on("console", lambda m, e=errors: m.type == "error"
                            and not m.text.startswith("Failed to load resource") and e.append(m.text))
                    page.on("pageerror", lambda x, e=errors: e.append(str(x)))
                    page.on("requestfailed", lambda r, f=failed: f.setdefault(r.url, r.failure or "failed"))
                    page.on("response", lambda r, f=failed: r.status >= 400 and f.update({r.url: str(r.status)}))
                    page.goto(to_url(base, root, path), wait_until="load")
                    page.evaluate("t => document.documentElement.setAttribute('data-theme', t)", theme)
                    page.wait_for_timeout(200)
                    overflow = page.evaluate(
                        "document.documentElement.scrollWidth - document.documentElement.clientWidth")
                    name = path.relative_to(root).with_suffix("").as_posix().replace("/", "__")
                    shot = out / f"{name}-{width}-{theme}.png"
                    page.screenshot(path=str(shot), full_page=True)
                    page.close()
                    report["shots"].append({"page": str(path.relative_to(root)), "width": width, "theme": theme,
                                            "shot": str(shot), "console_errors": errors,
                                            "failed_assets": [f"{why} {u.replace(base, '')}" for u, why in failed.items()],
                                            "horizontal_scroll_px": overflow})
                    problems += len(errors) + len(failed) + (overflow > 0)
                ctx.close()
        browser.close()
    server.shutdown()

    (out / "look.json").write_text(json.dumps(report, indent=2))
    for s in report["shots"]:
        issues = s["console_errors"] + s["failed_assets"]
        if s["horizontal_scroll_px"] > 0:
            issues.append(f"scrolls sideways by {s['horizontal_scroll_px']}px")
        if issues:
            print(f"FAIL {s['page']} @{s['width']} {s['theme']}: " + "; ".join(issues))
    for m in report["missing_links"]:
        print(f"FAIL link in {m['from']} goes to missing {m['to']}")
    for u in report["unreachable"]:
        print(f"FAIL {u} cannot be reached from {args.crawl}")
    print(f"{len(report['shots'])} screenshots in {out}; {problems} problem(s). Now open the screenshots.")
    sys.exit(1 if problems else 0)


if __name__ == "__main__":
    main()
