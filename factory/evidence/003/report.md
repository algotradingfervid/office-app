# Evidence — story 003

story: factory/stories/003-leave-contract.md
code-commit: 9aeee67ce0a4cbe99612a3c7310636cbeb7d0025
base: origin/main (5387e0e)
captured: 2026-09-25 10:08
check: "The leave collections exist with the design's fields, the default rules for all nine leave types are seeded, and a ledger entry can never be edited or deleted"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/auth	5.254s
ok  	officeapp/internal/core/clock	0.548s
ok  	officeapp/internal/forms/leave	3.915s
```

## Story check commands

```
$ go test -count=1 ./internal/forms/leave/   -> exit 0
ok  	officeapp/internal/forms/leave	3.115s
```

## Mutation (scripts/mutate.sh 003)

```
$ scripts/mutate.sh 003   -> PASS
mutation: base origin/main; changed files:
  internal/forms/leave/1790310300_leave.go
  internal/forms/leave/ledger_guard.go
  internal/forms/leave/rules.go
  internal/forms/leave/seed_rules.go
mutation: killed 1, lived 0, not covered 0
```

## Verdict

check: "The leave collections exist with the design's fields, the default rules for all nine leave types are seeded, and a ledger entry can never be edited or deleted"
verdict: PASS — review round 1 fixes in (ledger index includes leave_year, year-close test; per-field 0 meanings; table-driven half-day test; LeaveYear doc). make check green. Scripted mutation: 1 killed, 0 lived, 16 timed out under load ~31; a manual rerun, gremlins --timeout-coefficient 40, killed 17, lived 0, timed out 0
