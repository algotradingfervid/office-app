---
id: 011
title: I can see my requests and their status
tier: core
lane: full
kind: feature
status: building
needs: ["006"]
files: ["internal/core/approvals/my_requests_page.go", "internal/core/approvals/templates/my_requests.html"]
screen: my requests (/requests)
check: "A signed-in employee sees a list of only their own requests with form summary, dates submitted and status, newest first"
agent: claude-bg-011
started: 2026-09-25 10:50
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# I can see my requests and their status

## What
`GET /requests` for signed-in users: the employee's own requests (summary from the form's `Summary`, submitted date in IST, status, current approver name), newest first, each linking to `/requests/{id}`.
A nav link "My requests" is not part of this story (the layout belongs to web); link from the page title is enough for now.

## Check
```check
go test -count=1 -run 'MyRequests' ./internal/core/approvals/
```
ApiScenario tests: E001 sees their request, not E002's; signed out → redirect to /login.

## Out of scope
Request detail page (015), filters, paging.

## Constraints
Embed only this page's template. Authorisation on the server.
