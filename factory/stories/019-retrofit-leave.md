---
id: 019
title: Restyle the apply-for-leave screen to its wireframe
tier: usable
lane: core
kind: feature
status: blocked
needs: ["017", "014"]
files: ["internal/forms/leave/templates/apply.html", "internal/forms/leave/templates/apply_preview.html"]
screen: leave screens
check: "The leave screens match their approved wireframes at 390 and 1280 px"
blocked: waiting for wireframe approval
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# Restyle the apply-for-leave screen to its wireframe

## What
Restyle the Core-tier leave templates to their approved wireframes using the design-system components from 017. Markup and classes only; no behaviour changes.

## Check
The screens' existing journeys (stories 007, 011, 013, 014, 015) re-run by `scripts/prove.sh 019`, screenshots beside the approved wireframes.

## Out of scope
Handlers, services, new features.

## Constraints
Released by /wireframe after 017 merges.
