---
id: 006
title: A request can be submitted to a first approver and forwarded
tier: core
lane: full
kind: feature
status: merged
needs: ["002"]
files: ["internal/core/approvals/service_submit.go"]
screen: none: no UI
check: "Submitting creates a pending request with a submitted step for the chosen approver; the current approver can forward it to another active approver, at most 5 times"
agent: claude-bg-006
started: 2026-09-25 10:21
built: 2026-09-25 10:25
proved: 2026-09-25 10:27
reviewed: 
merged: 2026-09-25 10:46
review-rounds: 0
pr: 
---

# A request can be submitted to a first approver and forwarded

## What
Service functions in approvals (design spec `docs/superpowers/specs/2026-09-25-office-app-leave-design.md` §7 steps 1–2):
- `Submit`: in one transaction, creates the `requests` row (pending, current_approver), lets the form save its own data through a callback, runs the form's `Validate` (errors roll everything back; warnings are stored on the step), writes step `submitted` (seq 1).
- `Forward`: only the current approver; target must be active, `can_approve`, not the actor, not the requester; max 5 forwards; writes step `forwarded` with comment and to_user.
Refusals are `UserError`s with plain messages.

## Check
```check
go test -count=1 -run 'Submit|Forward' ./internal/core/approvals/
```
Tests with a fake form registered on a test app: happy path; first approver = requester refused; inactive or non-approver target refused; non-current approver refused; 6th forward refused; a Validate error leaves no rows behind; warnings stored.

## Out of scope
Approve/reject (010), notifications (Usable), pages.

## Constraints
Status moves through `Next` from story 002. Every read inside the transaction uses `txApp`.
Before calling `Next`, check the actor against the transition table's "who" (current approver / requester / HR admin, active) inside the transaction; one test per role that fails if the check is removed.
