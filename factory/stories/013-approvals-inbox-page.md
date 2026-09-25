---
id: 013
title: Approvers see what is waiting for them
tier: core
lane: full
kind: feature
status: building
needs: ["010"]
files: ["internal/core/approvals/inbox_page.go", "internal/core/approvals/templates/inbox.html"]
screen: pending my approval (/approvals)
check: "An approver sees only requests where they are the current approver, with requester, summary and waiting time; a non-approver is refused"
agent: claude-bg-013
started: 2026-09-25 15:42
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# Approvers see what is waiting for them

## What
`GET /approvals` for users with `can_approve`: pending and cancel_requested requests whose current approver is me, oldest first, with requester name, summary, submitted date and days waiting, each linking to `/requests/{id}`.

## Check
```check
go test -count=1 -run 'Inbox' ./internal/core/approvals/
```
ApiScenario tests: E002 sees a request assigned to E002, not one assigned to E003; E001 (not an approver) gets 403; signed out → /login.

## Out of scope
Actions (015), counts on home, filters.

## Constraints
Embed only this page's template.
Filter the inbox by status (pending, cancel_requested) as well as current approver: rejected requests keep `current_approver` (010 security review).
