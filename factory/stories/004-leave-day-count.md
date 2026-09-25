---
id: 004
title: A leave request's days are counted like the policy says
tier: core
lane: full
kind: feature
status: proving
needs: ["001", "003"]
files: ["internal/forms/leave/daycount.go"]
screen: none: no UI
check: "Days are counted per the design: weekly offs and holidays skipped, half-day sessions count 0.5, maternity counts calendar days, and a request starting or ending on a non-working day is refused"
agent: claude-bg-004
started: 2026-09-25 10:21
built: 2026-09-25 10:24
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# A leave request's days are counted like the policy says

## What
One function in the leave module computes a request's days from its rule's `count_mode`, the from/to dates and sessions, using `calendar.WorkingDays` (design spec `docs/superpowers/specs/2026-09-25-office-app-leave-design.md` §5.3 "Day count").

## Check
```check
go test -count=1 -run 'Days' ./internal/forms/leave/
```
Table-driven, fixed dates: 2026-10-12..13 full = 2; 2026-10-09..12 = 2 (Sat 10 off, Sun 11 off); 2026-10-12 first_half = 0.5; 2026-10-12 second_half..2026-10-13 first_half = 1;
maternity 2026-10-01..2026-10-31 calendar = 31; start on 2026-10-11 (Sun) refused; half-day session on a holiday refused; to before from refused.

## Out of scope
Policy limits (008, 009), UI.

## Constraints
User-facing refusals are `approvals.UserError` messages in plain English.
