---
id: 023
title: Probation is judged by the leave dates, not the day of applying
tier: core
lane: full
kind: feature
status: proving
needs: ["008"]
files: ["internal/forms/leave/policy_dates.go"]
screen: none: no UI
check: "An employee on probation can apply for Earned Leave that starts on or after their probation end date, and is refused for Earned Leave that starts before it"
agent: claude-bg-023
started: 2026-09-25 15:42
built: 2026-09-25 16:00
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# Probation is judged by the leave dates, not the day of applying

## What
Owner decision 2A (2026-09-25; design spec §5.3 check 1 updated): a leave request is probation leave when its
`from_date` is before the employee's `probation_end` (the confirmation date, the first day after probation).
Story 008 judged it by today's date. Change check 1 in `DateChecks` accordingly; empty `probation_end` = no probation.

## Check
```check
go test -count=1 -run 'DateChecks' ./internal/forms/leave/
```
Table cases (clock 2026-10-05, probation_end 2026-10-20): EL from 2026-11-10 allowed; EL from 2026-10-19 refused;
EL from 2026-10-20 allowed; CL during probation still allowed; empty probation_end allowed. Update 008's existing
probation cases in policy_dates_test.go to the new rule.

## Out of scope
Any other check in DateChecks.

## Constraints
Keep the refusal message plain; mention the confirmation date.
