---
id: 008
title: Leave dates are checked before submit
tier: core
lane: full
kind: feature
status: backlog
needs: ["002", "003", "021"]
files: ["internal/forms/leave/policy_dates.go"]
screen: none: no UI
check: "A leave request is refused when its type is inactive or barred in probation, it crosses 31 March, it overlaps the employee's own leave, or it is backdated too far; short notice is only a warning"
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# Leave dates are checked before submit

## What
Checks 1–5 of design spec `docs/superpowers/specs/2026-09-25-office-app-leave-design.md` §5.3 as one exported function in the leave module returning (warnings, error):
type active and allowed in probation; valid dates and no leave-year crossing; no overlap with own pending/approved/cancel_requested requests (half-day aware);
`from_date ≥ today − max_backdate_days`; `from_date < today + min_notice_days` → warning "short notice".

## Check
```check
go test -count=1 -run 'DateChecks' ./internal/forms/leave/
```
Table-driven with `clock.Fixed` 2026-10-05: each rule's pass and fail case; first_half vs second_half on the same date does not overlap; a cancelled request does not overlap.

## Out of scope
Limits and balance (009). Aggregation into `Validate` (012).

## Constraints
Messages are plain sentences an employee understands.
