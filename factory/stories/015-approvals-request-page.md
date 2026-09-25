---
id: 015
title: A request page shows its history and the approver's actions
tier: core
lane: full
kind: feature
status: backlog
needs: ["010", "012"]
files: ["internal/core/approvals/request_page.go", "internal/core/approvals/templates/request.html"]
screen: request detail (/requests/{id})
check: "The requester and approvers see every step with who, when and comment; the current approver can forward, approve as final or reject from the page"
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# A request page shows its history and the approver's actions

## What
`GET /requests/{id}`: summary, status, warnings, and the history timeline (submitted, forwarded, approved…) with actor, IST time and comment. Visible to the requester, anyone who acted on it, and the current approver; others get 404.
For the current approver: forms to Approve & forward (next approver + comment), Approve as final, Reject (comment required), posting to `/requests/{id}/forward|approve|reject`, which call the service and redirect back.

## Check
```journey
login E001
goto /leave/apply
select [name=leave_type] CL
fill [name=from_date] 2026-10-12
fill [name=to_date] 2026-10-13
select [name=first_approver] Ravi Kumar
click text=Submit
see Pending
login E002
goto /approvals
click text=CL
select [name=to_user] Meena Iyer
click text=Approve & forward
login E003
goto /approvals
click text=CL
click text=Approve as final
see Approved
see Ravi Kumar
shot history
```
```check
go test -count=1 -run 'RequestPage' ./internal/core/approvals/
```

## Out of scope
Cancel buttons (Usable). Leave-specific detail beyond `Summary`.

## Constraints
Embed only this page's template. Permission checks on the server for GET and every POST.
An unknown or inaccessible request id returns 404 with a plain page, never a raw error; every POST re-checks the actor server-side (006 security review).
