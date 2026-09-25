---
id: 002
title: Requests, approval steps and the status table exist
tier: core
lane: full
kind: contract
status: review
needs: []
files: ["internal/core/approvals/1790310200_requests.go", "internal/core/approvals/transitions.go", "internal/core/approvals/forms.go", "internal/core/approvals/history_guard.go"]
screen: none: no UI
check: "A request's status can only move along the design's transition table, and an approval step can never be edited or deleted"
agent: claude-bg-002
started: 2026-09-25 09:56
built: 2026-09-25 09:59
proved: 2026-09-25 10:00
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# Requests, approval steps and the status table exist

## What
The approvals contract every form builds on (design spec `docs/superpowers/specs/2026-09-25-office-app-leave-design.md` §4, §5.1, §7).
- Collections `requests` and `approval_steps` with the fields of §5.1 (skip `recorded_by_hr` behaviour; the field exists), API rules `nil`, unique (request, seq) on steps.
- `approval_steps` is append-only: `OnRecordUpdate`/`OnRecordDelete` hooks for that collection return an error (also blocks the dashboard).
- Contract (exported from the root package): status and action constants; the §7 transition table as data with
  `Next(from Status, a Action) (Status, bool)`; the form hooks interface
  `Form{ Validate(txApp core.App, c clock.Clock, req *core.Record) (warnings []string, err error); OnFinalApproved(txApp core.App, req *core.Record) error; OnCancelledAfterApproval(txApp core.App, req *core.Record) error; Summary(app core.App, req *core.Record) (string, error) }`,
  and `RegisterForm(app core.App, formType string, f Form)` / `FormFor(app core.App, formType string) (Form, bool)` (stored per app, like `home.AddCard`).
- A `Validate` error meant for the user is a `*UserError` (type in this package) whose message is shown as-is; any other error is internal.

## Check
```check
go test -count=1 ./internal/core/approvals/
```
Tests: every row of the §7 table allowed and a sample of forbidden moves refused; saving a step works, updating or deleting it fails; RegisterForm/FormFor round-trip on a test app.

## Out of scope
The service actions (stories 006, 010) and pages.

## Constraints
Roles in the table are documentation here; permission checks live in the service stories.
