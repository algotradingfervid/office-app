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
