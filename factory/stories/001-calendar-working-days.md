---
id: 001
title: The calendar knows weekly offs and holidays
tier: core
lane: core
kind: contract
status: ready
needs: []
files: ["internal/core/calendar/1790310100_calendar.go", "internal/core/calendar/workdays.go", "internal/core/calendar/seed.go"]
screen: none: no UI
check: "For 2026-27, Sundays, 2nd and 4th Saturdays and seeded holidays are non-working days; every other date is a working day"
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# The calendar knows weekly offs and holidays

## What
The shared calendar module gets its data and the one question every form asks: is this date a working day?
- Collections `holidays` (date text `YYYY-MM-DD` unique, name) and `settings` (key unique, value JSON), API rules `nil` (design spec `docs/superpowers/specs/2026-09-25-office-app-leave-design.md` §5.2).
- Default setting `weekly_off`: Sundays plus the 2nd and 4th Saturday of each month.
- Exported from the root package (this is the contract): `IsWorkingDay(app core.App, date string) (bool, error)` and
  `WorkingDays(app core.App, from, to string) ([]string, error)` (working dates in the inclusive range, in order).
- Demo seed: the national holidays of leave year 2026-27 (Republic Day 2027-01-26, Independence Day 2026-08-15, Gandhi Jayanti 2026-10-02) plus Diwali 2026-11-09.

## Check
```check
go test -count=1 ./internal/core/calendar/
```
Table-driven tests with fixed dates: 2026-10-10 (2nd Sat) off, 2026-10-17 (3rd Sat) working, 2026-10-24 (4th Sat) off, 2026-10-31 (5th Sat) working, 2026-10-11 (Sun) off, 2026-10-02 (holiday) off, 2026-10-12 working; `WorkingDays` over 2026-10-09..2026-10-13 returns 09, 12, 13.

## Out of scope
HR screens for holidays or settings (Usable). Half days and leave counting (story 004).

## Constraints
Dates are strings, never PocketBase `date` fields. Register the seeder with `seed.Add` via `addPart`.
