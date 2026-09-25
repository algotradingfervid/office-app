---
name: spike
description: "Prove or disprove the product's riskiest 'can we do it' assumption (data, accuracy, feasibility) with a throwaway script on real data and a measurable bar, before the spec is final. Use after /research when the brief's riskiest assumption is not demand, when the user says \"spike\", \"will this even work\", \"prove the core\", or invokes /spike."
argument-hint: "[the assumption, or empty to take it from factory/BRIEF.md]"
---

# /spike — Studio, station 3

Build the ugliest thing that answers one question with a number. It is not version one. Nothing here is reviewed, proved or merged into the product.

## Steps

0. Write the station marker so the guardrails let this station do its job: `echo spike > .claude/station`. Remove it when done: `rm -f .claude/station`.
1. Read the riskiest assumption in `factory/BRIEF.md` (or `$ARGUMENTS`) and `factory/research/feasibility.md`.
2. **Agree the bar with the human** before writing code: the question, the real data it runs on (e.g. 30 real documents the human has labelled or approved), the number that means go (e.g. "false positives ≤ 10%, recall ≥ 80%"), and a time budget (default: one day). Write them at the top of `factory/SPIKE.md`.
3. **Build the core loop only**, in `spikes/<name>/`: input → the hard step → output → the measurement. One script or notebook. No UI, no accounts, no tests beyond the measurement, no worktree.
4. **Measure** on the agreed data. Record the command and its output. Try at most three approaches; record the number for each.
5. **Write** `factory/SPIKE.md`: the bar, each approach with its number and command, what the numbers mean, and a verdict: **go** (bar met: the approach carries into `/define` and `/blueprint`), **change approach** (say which), or **stop** (the product as briefed does not work; say what would have to change).
6. Present the verdict and the numbers to the human. Set `spike: done` in `factory/STATE.md`.

## Gate

A number for the bar, produced by a command anyone can rerun, and the human's go / change / stop decision recorded in `SPIKE.md`.

## Rules

- The measurement data is the asset; the code is disposable. Keep the data and the scoring command: they become the Core tier's acceptance check (an eval) on the line.
- Stop at the time budget. An unanswered question is a finding; write it down and ask the human.
- Do not polish. If you are adding error handling, config or structure, you are building the product, not the spike.
