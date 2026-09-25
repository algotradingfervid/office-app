---
id: 022
title: Available days subtract leave that is still pending
tier: core
lane: full
kind: feature
status: backlog
needs: ["005", "021"]
files: ["internal/forms/leave/balance_available.go"]
screen: none: no UI
check: "Available days are the balance minus days of that type and leave year in the employee's pending requests; approved, rejected and cancelled requests hold nothing"
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# Available days subtract leave that is still pending

## What
Split from story 005 (it needs the leave_requests → requests link from 021). One exported function in the leave
module: `Available(app, employeeID, typeCode, leaveYear)` = `Balance(...)` (story 005) − days of that type and
leave year in the employee's `pending` requests (design spec §5.4). Also extend `Balances` output with available
if 005 returns a struct that has room for it — otherwise add a sibling function; do not edit 005's file.

## Check
```check
go test -count=1 -run 'Available' ./internal/forms/leave/
```
Tests: seeded E001 CL balance 4 and a pending 1-day CL → available 3, balance 4; an approved or rejected or cancelled request holds nothing; a pending EL does not affect CL; a pending request of 2025-26 does not affect 2026-27.

## Out of scope
Writing debits (012), pages.

## Constraints
Read through the `app` passed in (a `txApp` inside transactions). Ledger `days` are never 0.
