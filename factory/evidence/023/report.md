# Evidence — story 023

story: factory/stories/023-leave-probation-by-leave-date.md
code-commit: 6e6d5f2fc383cd3f6aef6a61b6ca4eca4b0991e0
base: origin/main (merge-base 1190dde; origin/main now 1190dde)
captured: 2026-09-25 15:44
check: "An employee on probation can apply for Earned Leave that starts on or after their probation end date, and is refused for Earned Leave that starts before it"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	2.077s
ok  	officeapp/internal/core/auth	2.340s
ok  	officeapp/internal/core/calendar	1.032s
ok  	officeapp/internal/core/clock	0.200s
ok  	officeapp/internal/forms/leave	0.990s
```

## Story check commands

```
$ go test -count=1 -run 'DateChecks' ./internal/forms/leave/   -> exit 0
ok  	officeapp/internal/forms/leave	0.695s
```

## Mutation (scripts/mutate.sh 023)

```
$ scripts/mutate.sh 023   -> PASS
mutation: base origin/main; changed files:
  internal/forms/leave/policy_dates.go
mutation: killed 2, lived 0, not covered 0, timed out 0
```

## Verdict

check: "An employee on probation can apply for Earned Leave that starts on or after their probation end date, and is refused for Earned Leave that starts before it"
verdict: PASS — Check 1 now uses from_date vs probation_end; story's table cases pass, make check green, 2/2 mutants killed
