---
id: 000
title: <observable change>
tier: core                # core | usable | launch | later
lane: core                # core | full  (full when the footprint hits a risk path in the CLAUDE.md product section)
kind: feature             # feature | contract (defines an interface others build on)
status: backlog
needs: []
files: []                 # exact files or narrow globs, ≤ 5 non-test files; this is an ownership claim
screen: none
check: "<a sentence a non-programmer could verify by looking>"
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# <title>

Keep this file under 60 lines. Say what must be true, not how to code it: no function signatures or file-by-file steps (a contract story is the exception: its interface is the deliverable).

## What
One paragraph: the user-visible change and where it lives.

## Check
How the check is run: the test, the command, or the screen and state to capture.

## Out of scope
What a tempted agent must not do in this story.

## Constraints
Decisions already made that apply: glossary terms, the contract this story builds on, the wireframe region, rules in ARCHITECTURE.md.
