# Evidence — story 022

story: factory/stories/022-leave-available.md
code-commit: 4e9b1536aefc6c1d4627879329fb5b5bb7eb4c40
base: origin/main (merge-base dfec94a; origin/main now dfec94a)
captured: 2026-09-25 10:59
check: "Available days are the balance minus days of that type and leave year in the employee's pending requests; approved, rejected and cancelled requests hold nothing"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	60.331s
ok  	officeapp/internal/core/auth	35.505s
ok  	officeapp/internal/core/calendar	18.437s
ok  	officeapp/internal/core/clock	6.418s
ok  	officeapp/internal/forms/leave	56.291s
```

## Story check commands

```
$ go test -count=1 -run 'Available' ./internal/forms/leave/   -> exit 0
ok  	officeapp/internal/forms/leave	24.860s
```

## Mutation (scripts/mutate.sh 022)

```
$ scripts/mutate.sh 022   -> PASS
mutation: base origin/main; changed files:
  internal/forms/leave/balance_available.go
mutation: killed 9, lived 0, not covered 0, timed out 0
```

## Verdict

check: "Available days are the balance minus days of that type and leave year in the employee's pending requests; approved, rejected and cancelled requests hold nothing"
verdict: PASS — pending 1-day CL gives available 3 with balance 4; approved/rejected/cancelled/cancel_requested, other type, year and employee hold nothing; 9/9 mutants killed
