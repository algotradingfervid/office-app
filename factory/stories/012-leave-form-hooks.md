---
id: 012
title: Leave plugs into approvals and debits on final approval
tier: core
lane: full
kind: feature
status: backlog
needs: ["008", "009", "005", "021"]
files: ["internal/forms/leave/form_hooks.go"]
screen: none: no UI
check: "The leave form is registered with approvals: its Validate runs the date and limit checks and recomputes days, final approval writes one debit per request, and its summary reads like 'CL · 12–13 Oct 2026 · 2 days'"
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# Leave plugs into approvals and debits on final approval

## What
The leave module registers form type `leave` with `approvals.RegisterForm` (design spec `docs/superpowers/specs/2026-09-25-office-app-leave-design.md` §4, §5.3 "Final approval"):
- `Validate`: loads the leave_request for the request, recomputes days (store the new value if a holiday changed it), runs date checks (008) and limit checks (009), returns warnings/errors.
- `OnFinalApproved`: one `debit` ledger entry (−days, leave year of from_date, linked to the request).
- `OnCancelledAfterApproval`: one `reversal` entry (+days). `Summary`: type code, date range, days.

## Check
```check
go test -count=1 -run 'LeaveForm' ./internal/forms/leave/
```
Integration test on a test app: submit a 2-day CL for E001 through `approvals.Submit`, approve final by E002 → CL balance 4 → 2, exactly one debit; a 3-day CL submit is refused with the policy message and leaves no rows.

## Out of scope
Pages. Year-close recompute on late changes (Usable).

## Constraints
Hooks use only the `txApp` they receive.
Test two submits of 1 CL each against 1 available: the second must be refused, because Validate reads `Available(txApp, ...)` and sees the first pending hold (022 review).
At final approval run the date checks (008) with a clock fixed at the request's `submitted_at` (IST), so backdate and notice are measured from submission (spec §5.3, owner decision 1A); overlap and limits use current data.
