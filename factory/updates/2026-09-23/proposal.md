# Kit update 2026-09-23 → v0.3.0: core first, small parallel stories, ceremony by risk

Status: **applied on branch `kit/update-2026-09-23`** (from `kit/update-2026-09-16b`, v0.2.1), 2026-09-23. Decisions below were answered by the founder the same day.

Source: the first full product run on the kit (SahiBid, `~/Documents/TenderFinder/shipshow`, 2026-09-16 to 2026-09-23), its `factory/lessons/LOG.md`, and the founder's review of the run on 2026-09-23: "taking too long to develop anything. I should be able to develop the core features fast and improve and build final product feature by feature", plus two additions of their own: plan stories and implementation plans for parallel execution, and keep stories small enough to build and test quickly.

## What went wrong on v0.2.1

The line was not slow per commit (1,349 commits, 38 stories merged in 7 days). It was slow to reach a product whose core works, because effort was spread evenly and the riskiest part came last.

| # | Finding | Evidence in SahiBid |
|---|---|---|
| 1 | The core idea was proven last | Stories 001–030 built sign-in, Telegram, digest, plan gate, second company, password reset, legal pages. Stories 031–039 then replaced the reading and matching core (037: "retires the criteria scoring"); the golden-set accuracy check arrived at 027/038. Work built on the old matcher was redone. |
| 2 | One ceremony level for everything | `DIALS.md` is `full` in every room and is product-wide. Name, brand, design system and hi-fi wireframes all ran before any code. Every story, from a legal page to the matcher, got a worktree, prove, a PR and a 5/5-from-all review loop. |
| 3 | Stories were specs of code, and too big | Median 6 files per story, max 31 (037). Story files average ~185 lines; 006 is 689, 014 is 576, with function signatures. `/slice` already says "under 300 lines, under 8 files", but nothing checked it. |
| 4 | Shared hotspot files limited parallel work | Across story footprints: `internal/api/api.go` in 15 stories, `migrations/*.sql` in 14, `cmd/sahibid/main.go` in 8, `internal/web/web.go` in 6. Wave 3: 015 and 017 both built `feed.go`; seven builders at once timed out `make check`. |
| 5 | Review loop long | Story files record review round 4 or later more than 20 times. Reviewer mutation found the most real gaps but was improvised each time. |
| 6 | Evidence heavy and repeatedly stale | 133 MB in `factory/evidence/`; ~100 one-off `scripts/journey-*` files; "evidence stale after review fixes" logged three times (029, 001, 012) and still `pending` in LOG.md. 016 merged with no report. |
| 7 | STATE.md too big to be a floor board | 74 KB, read first by every session and agent. |
| 8 | An agent wrote outside its worktree | 023 wrote 288 lines into main's checkout by absolute path; the guard only stops commits. |
| 9 | No timing data | Station time per story is not recorded; only 11 stories have a start and merge commit that pair up (median ~1.9 h start→merge), so where time goes is guessed, not measured. |

## Principles for v0.3.0

1. **Prove the core before building the product.** The riskiest assumption gets a throwaway spike with a number, before the spec is final.
2. **Build in tiers.** Core → Usable → Launch → Later. A tier starts when the one before works end to end.
3. **Ceremony follows risk, per story.** Full treatment for security, money and personal data; light treatment for the rest; none for spikes (a station, not a story lane).
4. **Small stories, planned in parallel waves.** Every story fits one sitting; stories in a wave own disjoint files; contracts come first.
5. **Modular monolith by default.** Parallel work comes from module boundaries and file ownership, not from separate services.
6. **Mechanisms over rules.** A lesson that recurs becomes a script or a hook, not another line of text.

Every change below names what it replaces. Removals are listed first where they exist.

## Changes

### A. Core first and tiers

| File | Change | Replaces | Risk |
|---|---|---|---|
| `.claude/skills/spike/SKILL.md` (new) | Throwaway proof of the riskiest assumption on real data, with a measurable bar agreed with the human (e.g. "false-positive rate ≤ X on a 30-bid golden set"). Code lives in `spikes/<name>/`, no worktree, no review, no evidence pack. Output `factory/SPIKE.md`: the bar, the number, the command that produced it, verdict (go / change approach / stop). Gate: a number, not an opinion. | The "disposable prototype" paragraph in `/brainstorm`, which only covered demand | medium |
| `/brainstorm` | Step 5 routes the riskiest assumption: demand → prototype as today; "can we do it" (data, accuracy, feasibility) → `/spike` after `/research`. | — | low |
| `/define`, `factory/templates/spec.md` | The ten actions are grouped into tiers: **Core** (the loop that proves value), **Usable** (one real user can use it), **Launch** (billing, legal, alerts, polish), **Later**. The human approves the tiers with the spec. | Flat list of ten actions | **high: changes what the spec gate approves** |
| `/slice` | Slices one tier at a time. The next tier is sliced only when the current tier's end-to-end check passes. | Slicing the whole spec at once (SahiBid: 27 stories up front) | medium |
| Station order (`CLAUDE.md`, `README.md`, `/shipshow`) | Two floors. **Prove it:** brainstorm → research → spike → define (H) → blueprint (stock styling) → tool-up → line runs Core. **Make it a product:** name-it (H) → brand (H) → design-system → wireframe (H), which can run *while* the line builds Core, then the line runs Usable and Launch against the approved design. One retrofit story restyles Core screens. | Nine upstairs stations in sequence before any code | **high: moves four human gates later** (the gates stay; their timing changes) |
| `/dials` | A new row, `design-timing: early \| after-core`. `early` keeps today's order (for products where the look is the risk); `after-core` is the default. | — | low |

### B. Ceremony by risk, per story

| File | Change | Replaces | Risk |
|---|---|---|---|
| `factory/templates/story.md`, `/slice` | New frontmatter `lane: core \| full` (throwaway measurements are `/spike`, not stories). `full` is set by path triggers that `/tool-up` writes into the product section (auth, payments, personal data, money, destructive migrations). Everything else is `core`. | One lane for every story | medium |
| `/ship` | **Gate becomes "no blocking findings"**, not "5/5 from every reviewer". Blocking = correctness, spec gap, security, data loss, untested guard. Non-blocking findings are recorded and not looped on. Maximum **two** fix rounds; a third goes to the human with the disagreement. `core` lane: `reviewer-code` only (+ `reviewer-design` if the story has a screen). `full` lane: plus every `product-review-*` whose trigger matches. | 5/5 from all, up to five loops | **high: changes the review loop** |
| `.claude/agents/reviewer-code.md` | Each finding is tagged `BLOCKING` or `NOTE`; the first line becomes `BLOCKING: n` alongside the score. | Score-only gate | medium |
| `scripts/mutate.sh` (written by `/tool-up` for the stack) | Runs mutation testing on the story's changed files only and lists surviving mutants; `/prove` runs it and puts the result in the report. The tool is chosen at `/tool-up` and checked for being current. | Reviewer mutation improvised each round (lesson 012) | medium |

### C. Small stories, parallel waves

| File | Change | Replaces | Risk |
|---|---|---|---|
| `/slice` size rule | One outcome, one check, **≤ 5 non-test files, ≤ ~300 changed lines**, buildable and provable by one fresh agent in one sitting. Over the limit → split along a file boundary. | Same rule (8 files) stated but not enforced | low |
| `/slice` story body | ≤ 60 lines: What, Check, Out of scope, Constraints. No function signatures or file-by-file instructions; the builder designs within the footprint. Exception: a **contract story**, whose deliverable *is* the interface. | 185-line average stories (max 689) | medium |
| `/slice` contract-first | A feature that more than one story will use starts with a contract story (types, interfaces, table, fake). The stories that build on it depend only on the contract, so they run side by side. | Implicit; SahiBid's 032 did this and it worked | low |
| `scripts/check-stories.sh` (new), `/slice` gate | Fails when a story is over the size limit, has no check, claims a whole directory with `*`, or claims a file another story in the same wave claims. | The unchecked gate sentence | low |
| `scripts/ready.sh` | Prints **waves**: the ready stories grouped so no two in a wave share a claimed file, plus a **hotspot report** (files claimed by ≥ 3 stories), which is fed back to `/blueprint`'s seam map. Recommended agents = min(wave width, merges the human can review per sitting). | Flat ready list; collision check only at `/isolate` | low |
| `/isolate`, orchestration notes | A story that follows another on shared files waits for it to merge; no "if it hasn't merged, you own the file" fallback. | Lesson 2026-09-22 (015/017), currently only in SahiBid's workflow scripts | low |
| `.claude/hooks/guard-worktree.sh` (new) | `PreToolUse` on Edit/Write: when the session's working directory is inside `.claude/worktrees/story-*`, block any `file_path` outside that worktree. | Lesson 2026-09-23 (023); an instruction in builder prompts | low |

### D. Architecture: modular monolith

| File | Change | Replaces | Risk |
|---|---|---|---|
| `/blueprint`, `factory/templates/architecture.md` | Default shape is a **modular monolith**: one repo, one binary, one deploy. Each feature is a module (folder) that owns its routes, migrations, queries and tests. Modules talk through small interfaces. **No central files to edit per story:** modules register their own routes and wiring; migrations are named by timestamp, not sequence. A separate process only for a workload reason (e.g. a background fetch/read worker, same codebase). Microservices only with a named reason (separate teams, separate scaling) recorded in ARCHITECTURE.md. | "Seams for parallel work" as a paragraph | medium |
| `/blueprint` step 4 | `make check` includes an import-boundary check, so one module cannot reach into another's internals. | Boundaries held by instruction | low |

### E. State, evidence, lessons

| File | Change | Replaces | Risk |
|---|---|---|---|
| `factory/templates/state.md`, `/shipshow` | STATE.md stays **≤ 60 lines**: current tier, stations, line counts, what's waiting on the human, open decisions. History and dated notes go to `factory/HISTORY.md`, which agents do not read by default. `ready.sh` warns when STATE.md passes 60 lines. | Unbounded STATE.md (SahiBid: 74 KB) | low |
| `/prove`, `scripts/prove.sh` | Screenshots only for stories with a screen; others prove with test and command output. `prove.sh` writes `report.md` itself, including the commit it proved. One reusable journey runner (written at `/tool-up`) reads per-story steps from the story file. Raw captures go to a git-ignored folder; git keeps `report.md` and the after-screenshots at the approved widths. | ~100 per-story journey scripts; 133 MB of evidence in git | medium |
| Merge step (`/ship` step 7) | Refuses to merge when `report.md` is missing, is still the template, or records a commit older than the branch head. It re-runs `prove.sh` instead of asking. | "Evidence stale" lesson ×3 still pending; 016 merged with no report | low |
| `/learn` | Runs once per **wave**, not per merge. A lesson that recurs after its rule was written becomes a script or a hook (no third rule). Pending LOG.md rows older than one wave are shown by `/shipshow`. | Per-merge /learn; rules that did not stop the repeat | low |

### F. The human's time, and measuring

| File | Change | Replaces | Risk |
|---|---|---|---|
| `/ship` hand-off, `/shipshow` | Merges are presented **per wave** as one review queue: one line per story (check, verdict, blocking count, evidence link). | One hand-off per story | low |
| Story frontmatter, `scripts/ready.sh --stats` | Record `started`, `built`, `proved`, `reviewed`, `merged` timestamps and `review-rounds`. `--stats` prints median time per station and rounds per lane, so the next update is based on measurements. | No timing data (finding 9) | low |

### G. Models (needs your decision; see below)

| File | Change | Replaces | Risk |
|---|---|---|---|
| `factory/MODELS.md`, `.claude/agents/researcher.md` | Two tiers only: Opus for judgement and code, Haiku for mechanical steps; researcher `model: sonnet` → `opus`. | Four tiers incl. Sonnet and Fable | low |

## Not changed

- Five human decisions stay five: spec, name, brand, wireframes, merge. Only their timing moves (A).
- Fresh context per story, writer never reviews its own work, evidence over claims, hooks over requests.
- The talk-first design process from v0.2.0/v0.2.1 stays as it is; it just runs on the second floor by default.

## Decisions (answered 2026-09-23)

1. **Station order (A):** founder asked for Claude's recommendation → `after-core` by default (design runs alongside Core, so it costs no wall-clock time; the talk-first design process is unchanged; `early` stays available when the look is the risk).
2. **Review gate (B):** yes.
3. **Models (G):** yes, kit default.
4. **v0.2.1:** yes, merge both into `main` together.

## Rolling it into SahiBid — on hold

The founder said on 2026-09-23 not to change TenderFinder now. When they ask, only the parts that help from here on:
- Move STATE.md history to `factory/HISTORY.md`; STATE.md back to one screen.
- Mark the matching work (038/039 and anything feeding accuracy) as the Core tier and finish it against the golden set before new features.
- Add `lane:` to open stories; run the new `ready.sh` hotspot report and turn `api.go`, `main.go` and sequential migrations into self-registration and timestamped migrations as one refactor story.
- Add the worktree guard hook and the stale-evidence merge check.
- Do not rewrite merged stories.

## Verification (when applying)

1. `scripts/check-stories.sh` and the new `ready.sh` run read-only against SahiBid's 39 stories: expect the size failures (e.g. 037) and the hotspots above to be reported.
2. `/slice` a new feature request with the new rules in a fresh subagent: every story ≤ 5 files and ≤ 60 lines, a contract story first, at least two stories in the first wave.
3. The tool-up trial (step 11) on the smallest story, `core` lane, end to end with the new `/ship` gate; record rounds and wall-clock via `--stats`.
4. The worktree guard blocks a write outside the worktree and allows one inside (both tested).
5. A `/spike` trial on a toy assumption produces SPIKE.md with a number and its command.
6. Frontmatter parses; no skill over 500 lines; CLAUDE.md under 200 lines; kit total lines not higher than v0.2.1 plus the new skill and scripts.
