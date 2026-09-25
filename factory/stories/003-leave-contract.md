---
id: 003
title: Leave types, rules, requests and the ledger exist
tier: core
lane: full
kind: contract
status: merged
needs: []
files: ["internal/forms/leave/1790310300_leave.go", "internal/forms/leave/rules.go", "internal/forms/leave/seed_rules.go", "internal/forms/leave/ledger_guard.go"]
screen: none: no UI
check: "The leave collections exist with the design's fields, the default rules for all nine leave types are seeded, and a ledger entry can never be edited or deleted"
agent: claude-bg-003
started: 2026-09-25 09:56
built: 2026-09-25 10:01
proved: 2026-09-25 10:09
reviewed: 
merged: 2026-09-25 10:21
review-rounds: 0
pr: 
---

# Leave types, rules, requests and the ledger exist

## What
The leave data contract (design spec `docs/superpowers/specs/2026-09-25-office-app-leave-design.md` §5.2, §5.4).
- Collections `leave_types`, `leave_rules`, `leave_requests` (relation to `requests`), `leave_ledger`, fields exactly as §5.2, API rules `nil`,
  partial unique index on `leave_ledger` (employee, leave_type, entry_type, period_key) where period_key != ''.
- `leave_ledger` append-only via update/delete hooks.
- Seed: the nine leave types (CO inactive) and one rule per type effective `2026-04-01` with the §5.2 default values.
- Contract (exported): `LeaveYear(date string) string` (`2026-10-12` → `2026-27`, `2027-03-31` → `2026-27`, `2027-04-01` → `2027-28`)
  and `RuleFor(app core.App, typeCode, date string) (*core.Record, error)` = the rule with the latest `effective_from` on or before date.

## Check
```check
go test -count=1 ./internal/forms/leave/
```
Tests: LeaveYear boundaries; RuleFor picks the right version when a second rule effective `2026-12-01` is added; ledger update/delete refused; duplicate (employee, type, credit, `2026-11`) refused while two entries with empty period_key are allowed.

## Out of scope
Balances (005), day counting (004), checks (008, 009).

## Constraints
Migration timestamp must sort after `1790310200` (requests). Days are numbers constrained to multiples of 0.5.
