# Kit changelog

## 0.3.0 — 2026-09-23

Core first, small parallel stories, ceremony by risk. Lessons from SahiBid's line (38 stories, 2026-09-17 to 2026-09-23) and the founder's review; details and risks in `2026-09-23/proposal.md`.

- New `/spike`: a throwaway proof of the riskiest "can we do it" assumption, with a number against an agreed bar. `/brainstorm` routes to it.
- `/define`: actions grouped into Core / Usable / Launch / Later tiers, approved with the spec. `/slice` cuts one tier at a time.
- `/dials`: `design-timing: after-core` (default) runs name, brand, design system and wireframes alongside the Core tier; `early` keeps the 0.2 order. `/blueprint` builds the skeleton on stock styling when after-core; `/tool-up` seeds a retrofit story.
- Stories: `tier`, `lane` (core / full), `kind` (feature / contract), station timestamps, review rounds; ≤ 5 non-test files, ≤ 60-line body, exact file claims. New `scripts/check-stories.sh`; `scripts/ready.sh` prints waves, hotspots, the STATE.md length warning and `--stats` (both backed by `scripts/stories.py`).
- `/ship`: reviewers by lane; findings tagged BLOCKING or NOTE; gate is no blocking findings in at most two fix rounds; merge refuses a missing, template or stale report; merges queued per wave.
- `/prove`: one journey runner, screenshots only for screens, `scripts/mutate.sh` on changed files, report pinned to the branch head, raw captures git-ignored.
- `/blueprint`: modular monolith by default, self-registering modules, timestamped migrations, import-boundary check in `make check`.
- New hook `guard-worktree.sh`: a story session cannot write outside its worktree. `/spike` may commit on main.
- `STATE.md` one screen, history in `HISTORY.md` (new template). `/learn` runs per wave and turns repeated lessons into hooks or scripts.
- `MODELS.md`: Opus and Haiku only; researcher agent on Opus.

## 0.2.1 — 2026-09-16

- `/brand`, `/design-system`, `/wireframe`: earlier remarks can shape the talk-step questions but do not count as answers; if the human is not available, stop and wait instead of drawing. `/wireframe` draws no other screens until the human has reacted to the first one. Gates say the talk happened at the station, before drawing.

## 0.2.0 — 2026-09-16

Design stations talk first, look, and get a blind review. Lessons from the first product run (four entries, 2026-09-16); details and risks in `2026-09-16/proposal.md`.

- `/brand`: direction talk (feel, references, colour appetite, logo approach, type) recorded in `factory/design/DIRECTION.md` and OK'd before any direction is drawn; look.py on previews; palette and mark talked through when presenting.
- `/design-system`: interface check-in before tokens; look.py on the style page; blind review; shown to the human.
- `/wireframe`: branded high-fidelity, responsive, clickable wireframes on the real design-system files; layout talk and one agreed screen first; look.py with link crawl; blind review; approved screenshots locked in `wireframes/approved/`.
- `/blueprint`, `/build`, `/prove`, `/ship`, `/tool-up`, `reviewer-design`: build and judge against the approved screenshots.
- New: `scripts/look.py`, `scripts/blind-copy.sh`, `.claude/agents/reviewer-blind.md`. Two house rules in `CLAUDE.md`.
