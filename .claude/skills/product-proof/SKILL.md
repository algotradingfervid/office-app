---
name: product-proof
description: "How to capture evidence for an Office App story: the prove script, the journey steps format for screen stories, check commands for non-screen stories, mutation testing, and how to read the report. Load in /prove, and when writing a story's Check section in /slice."
user-invocable: false
---

# Office App proof playbook

One command does the capture: `scripts/prove.sh <id>` (run inside the story worktree). It writes
`factory/evidence/<id>/report.md` with `code-commit`, and fills it with:

1. `make check` result (must be PASS).
2. The story's ` ```check ` commands (shell lines in the Check section), each with exit code and output tail.
3. Screen stories only: the story's ` ```journey ` block run by `scripts/journey.py` on this branch (`after/`)
   and on `origin/main`/`main` (`before/`), at **390 px and 1280 px, light theme**, each run on a fresh
   seeded database. Before the wireframes exist, judge structure and the check only (no styling review).
4. `scripts/mutate.sh <id>` — gremlins on changed lines. `LIVED` fails; `NOT COVERED` is listed for information
   (usually `if err != nil` returns). Install once: `go install github.com/go-gremlins/gremlins/cmd/gremlins@v0.6.0`.

The report starts with `verdict: PENDING`. /prove reads it, looks at every screenshot in `after/`, then
replaces that line with `verdict: PASS` or `verdict: GAP` and adds a `## Verdict` section: the check verbatim,
PASS/GAP, and one line per piece of evidence. Raw logs stay in `raw/` (git-ignored); commit `report.md` and `after/`.

## Writing the Check section of a story

Non-screen story — a check block of commands that exit 0 only when the check holds:

````
```check
go test -count=1 -run 'TestWorkingDays' ./internal/core/calendar/
```
````

Screen story — a journey block (steps below) **and** a check block for its tests:

````
```journey
login E001
goto /leave/apply
select [name=leave_type] CL
fill [name=from_date] 2026-10-12
fill [name=to_date] 2026-10-14
see at most 2 consecutive days
shot apply-too-long
```
````

Steps: `login <code>` (demo password), `goto <path>`, `fill <selector> <value>`, `select <selector> <value>`,
`click <selector>`, `see <text>`, `notsee <text>`, `wait <ms>`, `shot <name>`. Selectors are Playwright
selectors (`[name=x]`, `text=Submit`, `role=button[name="Approve as final"]`). Each width starts from the demo
data (E001–E004 and whatever the module seeders add), so a journey sets up what it needs through the UI or relies on seed data.

Never write a per-story script: extend `scripts/journey.py` with a new step verb in its own story if one is missing.

## Reading the evidence

- `make check` FAIL or any check command non-zero → GAP.
- Journey step FAIL on `after` → GAP. `before` failing on the new feature is expected; `before` failing on an
  old step (login, home) means the story broke something upstream of it — investigate.
- Screenshots: the page shows what the check says, nothing is cut off at 390 px, no raw error text.
- Any `LIVED` mutant in new logic → GAP: add the test that kills it, re-run `scripts/prove.sh`.
- After any code change the evidence is stale: rerun `scripts/prove.sh <id>`; never edit the report by hand
  except the verdict line and section.
