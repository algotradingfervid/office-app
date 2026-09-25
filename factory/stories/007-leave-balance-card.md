---
id: 007
title: Home shows my leave balances
tier: core
lane: core
kind: feature
status: proving
needs: ["005", "022"]
files: ["internal/forms/leave/home_card.go", "internal/forms/leave/templates/balance_card.html"]
screen: home (/)
check: "After signing in, E001 sees cards for CL 4, SL 8 and EL 10.5 with available days on the home page"
agent: claude-bg-007
started: 2026-09-25 15:42
built: 2026-09-25 16:05
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# Home shows my leave balances

## What
The leave module adds a home card (via `home.AddCard`) listing, for the signed-in employee and the current leave year, each active quota type with balance and available days.

## Check
```journey
login E001
see Casual Leave
see 10.5
shot home-balances
```
```check
go test -count=1 -run 'BalanceCard' ./internal/forms/leave/
```

## Out of scope
Pending-approval count, upcoming holidays (Usable). Any styling beyond stock Pico.

## Constraints
The card gets the leave year from `clock.System{}` via `LeaveYear`. The card parses its own template file (not the page layout).
