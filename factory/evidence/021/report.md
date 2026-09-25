# Evidence — story 021

story: factory/stories/021-leave-request-link.md
code-commit: 04354fb1d8c68d46bd783bc2003c41646add31a8
base: origin/main (6cebbb9)
captured: 2026-09-25 10:24
check: "Every leave request belongs to exactly one approval request, and a leave request cannot be saved without one"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	1.243s
ok  	officeapp/internal/core/auth	2.277s
ok  	officeapp/internal/core/calendar	1.848s
ok  	officeapp/internal/core/clock	1.649s
ok  	officeapp/internal/forms/leave	2.548s
```

## Story check commands

```
$ go test -count=1 -run 'LeaveRequestLink' ./internal/forms/leave/   -> exit 0
ok  	officeapp/internal/forms/leave	1.147s
```

## Mutation (scripts/mutate.sh 021)

```
$ scripts/mutate.sh 021   -> PASS
mutation: base origin/main; changed files:
  internal/forms/leave/1790310400_leave_request_link.go
mutation: killed 4, lived 0, not covered 0, timed out 0
```

## Verdict

check: "Every leave request belongs to exactly one approval request, and a leave request cannot be saved without one"
verdict: PASS — Tests show a leave request needs exactly one existing request (required, unique, no cascade) and ledger entries save with or without one; make check green; mutation 4/4 killed
