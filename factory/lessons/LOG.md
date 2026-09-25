# Lessons log

| date | story | lesson | classified as | kit change |
|---|---|---|---|---|
| 2026-09-25 | tool-up trial (001) | guard-main-branch.sh checked the session's branch, not the command's directory, and never matched `git -C <dir> commit`; a delegated agent could not commit in its worktree | hook bug (core) | fixed locally in guard-main-branch.sh (cd/-C aware); upstream via /shipshow-update |
| 2026-09-25 | tool-up trial (001) | /prove says write the verdict in report.md but guard-protected-files blocks factory/evidence/* for every station | core contradiction | product: `scripts/prove.sh <id> --verdict`; upstream: allow station=prove or script the verdict |
| 2026-09-25 | tool-up trial (001) | code-commit formula (`-- . ':!factory/evidence'`) counts the story-status commit, so evidence looks stale after /prove commits status: review | core contradiction | product-proof overrides the formula to also exclude factory/stories; upstream fix |
| 2026-09-25 | tool-up trial (001) | stop-check.sh writes /tmp/shipshow-check.log, shared by parallel worktrees | core race | pending upstream (use mktemp) |
| 2026-09-25 | tool-up trial (001) | isolate handoff assumes a human opens `claude -w`; no guidance for delegated agents; new-story.sh has no Go branch | core gap | product CLAUDE.md: `cd <worktree> &&` rule, Go needs no install; upstream |
| 2026-09-25 | tool-up | /tool-up left .claude/station and STATE `tool-up: pending` while the trial ran | station hygiene | remove marker before the trial; set STATE at the end |
| 2026-09-25 | tool-up trial 2 (001) | core /prove step 4 code-commit formula differs from product-proof's; guard-worktree is inactive for delegated agents and Bash-written files skip Edit/Write hooks | core gap | product: Write/Edit-only rule for delegated agents; upstream: path-based worktree guard, single code-commit formula |
| 2026-09-25 | wave 1-2 (003, 005) | twice /slice put stories in one wave where one needed data the other created (leave_requests→requests link; available→pending requests); fixed by new link stories 021, 022 | slicing rule (repeat) | pending /learn: /slice must list each story's reads of other stories' collections and put a link story after both |
| 2026-09-25 | wave 2 (004) | orchestrator `cd` into a worktree moved the shared session directory; a parallel builder's writes were then blocked by guard-worktree | orchestration | pending /learn: orchestrator uses absolute paths and `git -C`, never `cd`, while builders run |
| 2026-09-25 | wave 2 (004) | gremlins reports switch-case conditions as NOT COVERED although tests execute them | tool quirk | pending /learn: note in product-proof |
| 2026-09-25 | waves 1-2 | prove.sh records base origin/main (moving) instead of the branch merge-base | script | pending /learn: record `git merge-base` in report |
| 2026-09-25 | 004 | builder took the literal reading of an ambiguous spec sentence (maternity start/end rule); ledger reviewer caught it as BLOCKING | spec ambiguity | spec §5.3 clarified; reviewers working as intended |
