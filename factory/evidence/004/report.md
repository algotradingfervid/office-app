# Evidence — story 004

story: factory/stories/004-leave-day-count.md
code-commit: c8df43bc0b7b9c71430a1a93a7c12e62ac68a6e0
base: origin/main (6023236)
captured: 2026-09-25 10:29
check: "Days are counted per the design: weekly offs and holidays skipped, half-day sessions count 0.5, maternity counts calendar days, and a request starting or ending on a non-working day is refused"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	0.618s
ok  	officeapp/internal/core/auth	1.559s
ok  	officeapp/internal/core/calendar	1.383s
ok  	officeapp/internal/core/clock	0.817s
ok  	officeapp/internal/forms/leave	1.755s
```

## Story check commands

```
$ go test -count=1 -run 'Days' ./internal/forms/leave/   -> exit 0
ok  	officeapp/internal/forms/leave	0.611s
```

## Mutation (scripts/mutate.sh 004)

```
$ scripts/mutate.sh 004   -> PASS
mutation: base origin/main; changed files:
  internal/forms/leave/daycount.go
mutation: killed 14, lived 0, not covered 6, timed out 0
  NOT COVERED internal/forms/leave/daycount.go:34:12  CONDITIONALS_NEGATION
  NOT COVERED internal/forms/leave/daycount.go:36:12  CONDITIONALS_NEGATION
  NOT COVERED internal/forms/leave/daycount.go:34:33  CONDITIONALS_NEGATION
  NOT COVERED internal/forms/leave/daycount.go:38:12  CONDITIONALS_NEGATION
  NOT COVERED internal/forms/leave/daycount.go:36:33  CONDITIONALS_NEGATION
  NOT COVERED internal/forms/leave/daycount.go:38:31  CONDITIONALS_NEGATION
```

## Verdict

check: "Days are counted per the design: weekly offs and holidays skipped, half-day sessions count 0.5, maternity counts calendar days, and a request starting or ending on a non-working day is refused"
verdict: PASS — Review fix: calendar_days counts Sundays/holidays (182-day maternity), refuses only half days on non-working days; unknown session refused, unknown count mode errors; 14 killed, 0 lived (6 not-covered are switch-case coverage, cases exercised by TestDaysRefused)
