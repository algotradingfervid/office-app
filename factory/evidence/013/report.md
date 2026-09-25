# Evidence — story 013

story: factory/stories/013-approvals-inbox-page.md
code-commit: 61ae49a457a345714990464633845919aaafdd2f
base: origin/main (merge-base 1190dde; origin/main now 1190dde)
captured: 2026-09-25 15:45
check: "An approver sees only requests where they are the current approver, with requester, summary and waiting time; a non-approver is refused"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	1.688s
ok  	officeapp/internal/core/auth	1.515s
ok  	officeapp/internal/core/calendar	1.311s
ok  	officeapp/internal/core/clock	0.763s
ok  	officeapp/internal/forms/leave	1.815s
```

## Story check commands

```
$ go test -count=1 -run 'Inbox' ./internal/core/approvals/   -> exit 0
ok  	officeapp/internal/core/approvals	1.040s
```

## Journey (390 px and 1280 px, light; screen stories only)

```
after (61ae49a) -> PASS
PASS [390] login E002
PASS [390] goto /approvals
PASS [390] see Pending my approval
PASS [390] see Nothing is waiting for you.
PASS [390] shot inbox
PASS [1280] login E002
PASS [1280] goto /approvals
PASS [1280] see Pending my approval
PASS [1280] see Nothing is waiting for you.
PASS [1280] shot inbox
```
```
before (1190dde), expected to fail where the feature is new:
PASS [390] login E002
PASS [390] goto /approvals
FAIL [390] see Pending my approval  -> Locator.wait_for: Timeout 5000ms exceeded.
FAIL [390] console  -> Failed to load resource: the server responded with a status of 404 (Not Found)
PASS [1280] login E002
PASS [1280] goto /approvals
FAIL [1280] see Pending my approval  -> Locator.wait_for: Timeout 5000ms exceeded.
FAIL [1280] console  -> Failed to load resource: the server responded with a status of 404 (Not Found)
```

Screenshots: `after/` (this branch), `before/` (origin/main).

## Mutation (scripts/mutate.sh 013)

```
$ scripts/mutate.sh 013   -> PASS
mutation: base origin/main; changed files:
  internal/core/approvals/inbox_page.go
mutation: killed 6, lived 0, not covered 0, timed out 0
```

## Verdict

check: "An approver sees only requests where they are the current approver, with requester, summary and waiting time; a non-approver is refused"
verdict: PASS — E002 sees only pending/cancel_requested requests assigned to them, oldest first, with requester, summary, date and waiting days; E001 gets 403; signed out goes to /login; empty state renders at 390/1280; 6/6 mutants killed
