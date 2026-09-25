---
id: 005
title: Balances come from the ledger
tier: core
lane: full
kind: feature
status: review
needs: ["003"]
files: ["internal/forms/leave/balance.go", "internal/forms/leave/seed_balances.go"]
screen: none: no UI
check: "An employee's balance is the sum of their ledger entries for a type and leave year"
agent: claude-bg-005
started: 2026-09-25 10:21
built: 2026-09-25 10:24
proved: 2026-09-25 10:25
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# Balances come from the ledger

## What
Balance queries for the leave module (design spec `docs/superpowers/specs/2026-09-25-office-app-leave-design.md` §5.4): `Balance(app, employeeID, typeCode, leaveYear)`,
plus `Balances(app, employeeID, leaveYear)` for all active types with a quota (used by the home card), returned as a slice of a small exported struct (code, name, balance) so 022 can add available beside it.
Demo seed: `opening` entries for 2026-27 for E001–E004: CL 4, SL 8, EL 10.5.
Scope change (floor manager): `Available` moved to story 022, which needs 021's leave_requests → requests link.

## Check
```check
go test -count=1 -run 'Balance' ./internal/forms/leave/
```
Tests: seeded E001 CL = 4, SL = 8, EL = 10.5; after a −2 debit CL = 2; halves sum exactly; entries of 2025-26, of another type or of another employee do not count; `Balances` lists the active quota types only.

## Out of scope
Available days (022). Credits, year close, jobs (Usable). Writing debits (012).

## Constraints
Read through the `app` passed in (a `txApp` inside transactions). Sums of halves are exact; no rounding.
Rule fields use 0 two ways (see 003's migration comment): 0 = none/no limit for days_per_year, max_consecutive_days, yearly_cap, max_times_per_employment, attachment_after_days, min_notice_days; 0 = a real zero for max_backdate_days and carry_forward_cap. Ledger `days` cannot be 0 (PocketBase required number): never write a zero-day entry.
