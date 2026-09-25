# Evidence — story 008

story: factory/stories/008-leave-date-checks.md
code-commit: 31d5ccf0d63b686257d97ffe907635e1e599bdd1
base: origin/main (merge-base dfec94a; origin/main now dfec94a)
captured: 2026-09-25 11:03
check: "A leave request is refused when its type is inactive or barred in probation, it crosses 31 March, it overlaps the employee's own leave, or it is backdated too far; short notice is only a warning"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	41.196s
ok  	officeapp/internal/core/auth	44.864s
ok  	officeapp/internal/core/calendar	24.780s
ok  	officeapp/internal/core/clock	2.560s
ok  	officeapp/internal/forms/leave	49.549s
```

## Story check commands

```
$ go test -count=1 -run 'DateChecks' ./internal/forms/leave/   -> exit 0
ok  	officeapp/internal/forms/leave	20.541s
```

## Mutation (scripts/mutate.sh 008)

```
$ scripts/mutate.sh 008   -> PASS
mutation: base origin/main; changed files:
  internal/forms/leave/policy_dates.go
mutation: killed 35, lived 0, not covered 0, timed out 0
```

## Verdict

check: "A leave request is refused when its type is inactive or barred in probation, it crosses 31 March, it overlaps the employee's own leave, or it is backdated too far; short notice is only a warning"
verdict: PASS — Every check 1-5 has a pass and a fail case (half-day and cancelled no-overlap included); make check green; 35/35 mutants killed
