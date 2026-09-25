---
id: 005
title: Balances and available days come from the ledger
tier: core
lane: full
kind: feature
status: backlog
needs: ["003"]
files: ["internal/forms/leave/balance.go", "internal/forms/leave/seed_balances.go"]
screen: none: no UI
check: "An employee's balance is the sum of their ledger entries for a type and leave year, and available is that minus days held by pending requests"
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# Balances and available days come from the ledger

## What
Balance queries for the leave module (design spec `docs/superpowers/specs/2026-09-25-office-app-leave-design.md` §5.4): `Balance(app, employeeID, typeCode, leaveYear)` and `Available(...)` (balance − days of that type and leave year in `pending` requests),
plus `Balances(app, employeeID, leaveYear)` for all active types with a quota (used by the home card).
Demo seed: `opening` entries for 2026-27 for E001–E004: CL 4, SL 8, EL 10.5.

## Check
```check
go test -count=1 -run 'Balance|Available' ./internal/forms/leave/
```
Tests: seeded E001 CL = 4; after a −2 debit = 2; a pending 1-day CL request makes available 1 while balance stays 2; entries of 2025-26 do not count in 2026-27; approved/rejected requests do not hold days.

## Out of scope
Credits, year close, jobs (Usable). Writing debits (012).

## Constraints
Read through the `app` passed in (a `txApp` inside transactions). Sums of halves are exact; no rounding.
Rule fields use 0 two ways (see 003's migration comment): 0 = none/no limit for days_per_year, max_consecutive_days, yearly_cap, max_times_per_employment, attachment_after_days, min_notice_days; 0 = a real zero for max_backdate_days and carry_forward_cap. Ledger `days` cannot be 0 (PocketBase required number): never write a zero-day entry.
