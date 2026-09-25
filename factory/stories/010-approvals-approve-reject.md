---
id: 010
title: An approver can approve as final or reject
tier: core
lane: full
kind: feature
status: backlog
needs: ["006"]
files: ["internal/core/approvals/service_decide.go"]
screen: none: no UI
check: "The current approver can approve as final, which re-runs the form's checks and calls its approval hook exactly once, or reject with a required comment"
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# An approver can approve as final or reject

## What
Service functions in approvals (design spec `docs/superpowers/specs/2026-09-25-office-app-leave-design.md` §7 steps 3–4):
- `ApproveFinal`: only current approver; re-runs `Validate` inside the transaction; on success status `approved`, `final_approver` set, `current_approver` cleared, form's `OnFinalApproved` called, step `approved_final`.
- `Reject`: only current approver; comment required; status `rejected`; step `rejected`.
Status is re-read inside the transaction, so a second approve on the same request is refused.

## Check
```check
go test -count=1 -run 'ApproveFinal|Reject' ./internal/core/approvals/
```
Tests with a fake form counting hook calls: approve calls OnFinalApproved once; approving twice → second refused, hook count stays 1; a Validate error on approval leaves status pending; reject without comment refused; non-current approver refused.

## Out of scope
Cancellations, reassign (Usable). Pages.

## Constraints
Every read inside the transaction uses `txApp`.
