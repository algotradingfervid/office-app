---
id: 021
title: Leave requests are linked to their approval request
tier: core
lane: core
kind: contract
status: merged
needs: ["002", "003"]
files: ["internal/forms/leave/1790310400_leave_request_link.go"]
screen: none: no UI
check: "Every leave request belongs to exactly one approval request, and a leave request cannot be saved without one"
agent: claude-bg-021
started: 2026-09-25 10:21
built: 2026-09-25 10:24
proved: 2026-09-25 10:24
reviewed: 2026-09-25
merged: 2026-09-25 10:46
review-rounds: 1
pr: https://github.com/algotradingfervid/office-app/pull/4
---

# Leave requests are linked to their approval request

## What
Stories 002 and 003 were built in parallel, so `leave_requests` has no link to `requests` yet. A new migration
adds `leave_requests.request`: required, single relation to `requests`, cascade delete off, with a unique index
(one leave request per request), and `leave_ledger.request`: optional single relation to `requests` — design spec §5.2.

## Check
```check
go test -count=1 -run 'LeaveRequestLink' ./internal/forms/leave/
```
Tests: saving a leave_request without `request` fails; with a requests row it saves; a second leave_request for the same request fails; a ledger entry saves with and without `request`.

## Out of scope
Anything else on leave_requests.

## Constraints
Additive migration only; timestamp 1790310400 sorts after 002's and 003's migrations.
