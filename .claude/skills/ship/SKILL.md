---
name: ship
description: "Open a pull request for a proved story and run the review loop with fresh-context reviewers chosen by the story's lane, until no reviewer reports a blocking finding (at most two fix rounds), then queue it for the human to merge with the rest of its wave. Use after /prove, when the user says \"ship it\", \"PR\", \"review\", or invokes /ship <story-id>."
argument-hint: "<story-id>"
---

# /ship — the line, station 4 (human gate)

The writer never reviews its own work. Reviewers are fresh agents. The human reviews evidence, not diffs, and merges a wave at a time.

## Steps

1. Confirm `factory/evidence/$0/report.md` says PASS and its `code-commit` equals `git log -1 --format=%H -- . ':!factory/evidence'` (the last code change on the branch). If not, run `/prove $0` first.
2. **Open the PR** (`gh pr create`) from the story branch to `main`, using `factory/templates/pr.md`: what changed (three lines), how it was tested, the evidence, and risks. Title: `story $0: <title>`.
3. **Review.** Launch, in parallel and each with a clean context, the reviewers the story's `lane` calls for:
   - `core`: `reviewer-code`, plus `reviewer-design` if the story has a screen with an approved wireframe.
   - `full`: the above, plus every `.claude/agents/product-review-*.md` whose trigger matches the story's files (security for auth, payments, personal data; performance where a budget applies; platform for mobile or desktop).
   - If Greptile (or another external reviewer) is configured in the product section of CLAUDE.md, request its review too and treat its findings the same way.
   Each reviewer tags every finding `BLOCKING` or `NOTE` and gives a score.
4. **Fix the blocking findings only**: correctness, spec gap, security, data loss, an untested guard, a surviving mutant in new logic, or a screen that departs from its approved wireframe. Record NOTEs in the PR and do not loop on them. Chasing every note is how over-engineering gets in.
5. After fixes, run `scripts/prove.sh $0` and `scripts/mutate.sh $0` again (the evidence must describe the code that merges), push, and re-run only the reviewers that had blocking findings. **At most two fix rounds.** If a blocking finding is still open after the second, set `status: blocked` and summarise the disagreement for the human. Record `review-rounds` and `reviewed: <timestamp>` in the story file.
6. **Queue for merge.** Post a final PR comment: blocking count (zero), scores, notes, evidence summary, "Ready for merge". Set `status: review` with `pr: <url>`. The human merges once per wave: `/shipshow` shows every story in review as one queue, one line each (check, verdict, notes, evidence link).
7. **Before merge** (the human or the merge step): refuse if `report.md` is missing, is still the template, or its `code-commit` is not the branch's last code change; re-run `scripts/prove.sh $0` instead of asking.
8. **After the human merges** (detect via `gh pr view --json state` when `/shipshow` or `/learn` runs): set `status: merged` and `merged: <timestamp>`, and run `scripts/finish-story.sh $0` (deletes the worktree and branch, promotes stories whose `needs` are now met to `ready`). When the line drains (no story in building, proving or review) or five stories have merged since the last `/learn`, whichever comes first, run `/learn` over the stories merged since the last LOG.md entry.

## Gate

No open blocking finding, evidence current to the branch's last code change, PR opened. Then the human merges. The agent never merges.

## Rules

- `--force-with-lease` only, never `--force`.
- Regenerate lockfile and generated-code conflicts; never hand-edit them.
- A reviewer that always finds something is calibrated wrong; "nothing blocking, 5/5" is a valid result.
- Before any merge-path action, check the story's `status` and branch: after an outage or a resumed run, a story may already be merged.
