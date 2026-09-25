---
id: 009
title: Leave limits and balance are checked before submit
tier: core
lane: full
kind: feature
status: backlog
needs: ["002", "004", "005", "021"]
files: ["internal/forms/leave/policy_limits.go"]
screen: none: no UI
check: "A leave request is refused when it exceeds the consecutive-day limit (adjacent requests included), the available balance, the yearly cap, or the once-per-employment limit, or lacks a required attachment; LOP while paid leave remains is a warning"
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# Leave limits and balance are checked before submit

## What
Checks 6–11 of design spec `docs/superpowers/specs/2026-09-25-office-app-leave-design.md` §5.3 as one exported function in the leave module returning (warnings, error). Adjacent same-type requests separated only by non-working days are added together for the consecutive limit.

## Check
```check
go test -count=1 -run 'LimitChecks' ./internal/forms/leave/
```
Table-driven, fixed clock: 3 CL days refused (max 2); 2 CL Fri + 1 CL Mon after an existing Fri request refused (adjacency); 5 CL with CL available 4 refused;
LOP beyond yearly cap 15 refused; second marriage leave refused; SL 3 days without attachment refused, with attachment allowed; LOP while CL available → warning.

## Out of scope
Date checks (008), UI.

## Constraints
Days come from story 004's function; balance/available from 005.
