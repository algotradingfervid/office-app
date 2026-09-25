# Evidence — story 004

story: factory/stories/004-leave-day-count.md
code-commit: 39ce8aaf88b3ba875c2237c6be8b41f40d949948
base: origin/main (6cebbb9)
captured: 2026-09-25 10:24
check: "Days are counted per the design: weekly offs and holidays skipped, half-day sessions count 0.5, maternity counts calendar days, and a request starting or ending on a non-working day is refused"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	2.593s
ok  	officeapp/internal/core/auth	3.542s
ok  	officeapp/internal/core/calendar	3.550s
ok  	officeapp/internal/core/clock	2.877s
ok  	officeapp/internal/forms/leave	3.207s
```

## Story check commands

```
$ go test -count=1 -run 'Days' ./internal/forms/leave/   -> exit 0
ok  	officeapp/internal/forms/leave	2.289s
```

## Mutation (scripts/mutate.sh 004)

```
$ scripts/mutate.sh 004   -> PASS
mutation: base origin/main; changed files:
  internal/forms/leave/daycount.go
mutation: killed 13, lived 0, not covered 6, timed out 0
  NOT COVERED internal/forms/leave/daycount.go:28:12  CONDITIONALS_NEGATION
  NOT COVERED internal/forms/leave/daycount.go:28:33  CONDITIONALS_NEGATION
  NOT COVERED internal/forms/leave/daycount.go:30:12  CONDITIONALS_NEGATION
  NOT COVERED internal/forms/leave/daycount.go:30:33  CONDITIONALS_NEGATION
  NOT COVERED internal/forms/leave/daycount.go:32:12  CONDITIONALS_NEGATION
  NOT COVERED internal/forms/leave/daycount.go:32:31  CONDITIONALS_NEGATION
```

## Verdict

check: "Days are counted per the design: weekly offs and holidays skipped, half-day sessions count 0.5, maternity counts calendar days, and a request starting or ending on a non-working day is refused"
verdict: PASS — Days counts working/calendar days with half sessions and refuses non-working start/end, reversed range and bad session pairs; check green, 13 mutants killed, 0 lived (6 not-covered are gremlins switch-case coverage, cases are exercised by TestDaysRefused)
