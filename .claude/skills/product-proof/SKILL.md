---
name: product-proof
description: "How to capture evidence for an Office App story: the prove script, the journey steps format for screen stories, check commands for non-screen stories, mutation testing, and how to read the report. Load in /prove, and when writing a story's Check section in /slice."
user-invocable: false
---

# Office App proof playbook

**Commit the story's code first**; both scripts refuse uncommitted changes. Then, inside the story worktree:

```
scripts/prove.sh <id>                                  # capture: writes factory/evidence/<id>/report.md
scripts/prove.sh <id> --verdict PASS "one-line reason"  # after reading it (or GAP "what is missing")
```

The capture runs, in order: `make check`; the story's ` ```check ` commands (each must exit 0); for screen
stories only, the ` ```journey ` block via `scripts/journey.py` on this branch (`after/`) and on `main`
(`before/`) at **390 px and 1280 px, light theme**, each on a fresh seeded database; and `scripts/mutate.sh`
(gremlins on committed changed lines — `LIVED` fails, `NOT COVERED` is information). It already runs mutation:
do not run `scripts/mutate.sh` again separately (the /prove skill's step 3 is satisfied by this).
Install gremlins once: `go install github.com/go-gremlins/gremlins/cmd/gremlins@v0.6.0`.

The report is evidence: never edit it with Edit/Write (a hook blocks `factory/evidence/*`). The only change
after capture is `--verdict`, which refuses a stale report and also sets the story's `status` (PASS → `review` + `proved`, GAP → `building`). `code-commit` is the last commit touching anything
except `factory/evidence/` and `factory/stories/`; wherever the core skills compare `code-commit`
(/prove step 4, /ship steps 1 and 7), use `git log -1 --format=%H -- . ':!factory/evidence' ':!factory/stories'`,
so the story-status commit does not make the evidence look stale.
Look at every screenshot in `after/` before the verdict. Commit `report.md` and `after/` together with the
story's `status: review`; raw logs stay in `raw/` (git-ignored). Before wireframes exist, judge structure and the check only.

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
- After any code change the evidence is stale: rerun `scripts/prove.sh <id>`; never edit the report by hand.
