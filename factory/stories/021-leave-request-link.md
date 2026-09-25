---
id: 021
title: Leave requests are linked to their approval request
tier: core
lane: core
kind: contract
status: backlog
needs: ["002", "003"]
files: ["internal/forms/leave/1790310400_leave_request_link.go"]
screen: none: no UI
check: "Every leave request belongs to exactly one approval request, and a leave request cannot be saved without one"
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# Leave requests are linked to their approval request

## What
Stories 002 and 003 were built in parallel, so `leave_requests` has no link to `requests` yet. A new migration
adds field `request`: required, single relation to `requests`, cascade delete off, with a unique index
(one leave request per request) — design spec §5.2.

## Check
```check
go test -count=1 -run 'LeaveRequestLink' ./internal/forms/leave/
```
Tests: saving a leave_request without `request` fails; with a requests row it saves; a second leave_request for the same request fails.

## Out of scope
Anything else on leave_requests.

## Constraints
Additive migration only; timestamp 1790310400 sorts after 002's and 003's migrations.
