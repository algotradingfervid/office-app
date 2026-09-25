---
name: prove
description: "Capture evidence a non-programmer can judge for a story (screenshots beside the approved wireframe for screens; command output and numbers otherwise), run mutation testing on the changed code, write the report pinned to the branch head, and send the story back to /build if there is a gap. Use after /build, when the user says \"prove it\", \"show me\", \"evidence\", or invokes /prove <story-id>."
argument-hint: "<story-id>"
---

# /prove — the line, station 3

"Looks done" is not a signal. Evidence is. Evidence describes one commit; when the code changes, the evidence is recaptured.

## Steps

1. Load `product-proof` (the proof playbook written by `/tool-up`). It knows how this product is started and captured on each platform.
2. **Capture** with `scripts/prove.sh $0`. It runs the story's check (the steps written in the story's Check section, through the product's one journey runner; no per-story journey scripts), captures before (`origin/main`) and after (the story worktree), and writes `factory/evidence/$0/report.md` itself, including the `code-commit` it proved.
   - **Screen stories** capture screenshots at the widths and themes of `factory/design/wireframes/approved/`, placed beside the approved screenshots, with a list of every visible difference. Differences caused by real data are fine; different components, layout, colours, spacing or navigation are a GAP. On the Core tier before wireframes exist (`design-timing: after-core`), capture one width and theme and judge structure and the check only.
   - **Stories with no screen** prove with the test output, the command output and the numbers. No screenshots.
3. **Mutate.** Run `scripts/mutate.sh $0` on the story's changed files. A surviving mutant in new logic is a GAP: add the test that kills it. Put the result in the report.
4. **Verdict** in `report.md`: the acceptance check verbatim and one word, PASS or GAP; the evidence; the mutation result; the exact commands and their output; `code-commit: <sha>`, the last commit that changed anything outside `factory/evidence/` (`git log -1 --format=%H -- . ':!factory/evidence'`).
5. If GAP: list the gaps as concrete tasks, set `status: building`, and go back to `/build $0` yourself. Do not ask the human to look at a gap you can see. Maximum three prove→build loops; on the fourth, set `status: blocked` with the reason and stop.
6. If PASS: set `status: review` and `proved: <timestamp>`, commit `report.md` and the after-screenshots (raw captures, videos and logs stay in the git-ignored `factory/evidence/*/raw/`), then run `git ls-files factory/evidence/$0/report.md` to confirm the report is committed. Print `Next: /ship $0`.

## Gate

`report.md` is committed, says PASS, records the current `code-commit`, includes the mutation result, and every claim in it is backed by an artifact or a command's output.

## Rules

- For work with no visible surface (services, performance), the evidence is numbers and logs, with the command that produced them. A number without its command is a claim.
- Never edit the evidence. If it is wrong or its `code-commit` is older than the code, recapture.
