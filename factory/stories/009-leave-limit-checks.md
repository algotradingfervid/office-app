---
id: 009
title: Leave limits and balance are checked before submit
tier: core
lane: full
kind: feature
status: backlog
needs: ["002", "004", "005", "021", "022"]
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
Rule fields use 0 two ways (see 003's migration comment): 0 = none/no limit for days_per_year, max_consecutive_days, yearly_cap, max_times_per_employment, attachment_after_days, min_notice_days; 0 = a real zero for max_backdate_days and carry_forward_cap. Ledger `days` cannot be 0 (PocketBase required number): never write a zero-day entry.
`Balance` (005) returns 0 for an unknown type code: validate the leave type exists and is active before using balances.
The LOP yearly cap sums approved + pending LOP days itself; `AvailableBalances` lists quota types only (022 review).
