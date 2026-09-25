---
id: 016
title: The Core tier's end-to-end check runs with one command
tier: core
lane: core
kind: feature
status: backlog
needs: ["007", "011", "013", "014", "015"]
files: ["internal/e2e/core_test.go", "Makefile"]
screen: none: no UI
check: "`make e2e-core` on a fresh database proves: apply 2 CL, forward A→B, approve final → balance 4→2 with 3 history steps; 3-day CL refused; double approve deducts once; two concurrent submits exceeding balance → exactly one succeeds"
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# The Core tier's end-to-end check runs with one command

## What
An integration test package driving the real routes (ApiScenario / httptest with session cookies) through the Core value loop, and a `make e2e-core` target that runs it. This is the Core tier's gate in `factory/SPEC.md`.

## Check
```check
make e2e-core
```
Also run the journey of story 015 via `scripts/prove.sh 016` is not needed; this story has no screen.

## Out of scope
New features. If a step fails, write the gap to `factory/lessons/` and stop: the fix belongs to the owning story's module in a new story.

## Constraints
Only `internal/e2e/` and the one Makefile target. Concurrency: two goroutines submitting 3 CL each against 4 available → exactly one succeeds.
