---
id: 014
title: An employee can apply for leave with a live day count
tier: core
lane: core
kind: feature
status: backlog
needs: ["012", "006"]
files: ["internal/forms/leave/apply_page.go", "internal/forms/leave/templates/apply.html", "internal/forms/leave/templates/apply_preview.html"]
screen: apply for leave (/leave/apply)
check: "Picking dates shows working days and balance after without a reload; a 3-day CL shows the 2-day limit and cannot be submitted; a valid request appears as Pending"
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# An employee can apply for leave with a live day count

## What
`GET /leave/apply` (signed in): leave type (active types), from/to dates with sessions, reason, first approver (active `can_approve` users except me), submit.
Changing any field posts to a preview endpoint (htmx) that returns days, available and balance-after, and any errors/warnings from the same checks as submit.
`POST /leave/apply` calls `approvals.Submit` with a callback that saves the leave_request; on success redirect to `/requests/{id}`; on refusal re-render with the message.

## Check
```journey
login E001
goto /leave/apply
select [name=leave_type] CL
fill [name=from_date] 2026-10-12
fill [name=to_date] 2026-10-14
see at most 2 consecutive days
shot apply-too-long
fill [name=to_date] 2026-10-13
see 2 days
select [name=first_approver] Ravi Kumar
click text=Submit
see Pending
shot submitted
```
```check
go test -count=1 -run 'ApplyPage' ./internal/forms/leave/
```

## Out of scope
Attachments upload UI (Usable), styling beyond stock Pico.

## Constraints
The preview runs the same Go checks as submit; no client-side rules. No inline JS; htmx attributes only.
The requester is always the signed-in user (`e.Auth`), never a form field: include a test that a posted requester id is ignored (006 security review).
