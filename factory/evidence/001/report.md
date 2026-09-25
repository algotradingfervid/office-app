# Evidence — story 001

story: factory/stories/001-calendar-working-days.md
code-commit: af4cfda7ce7ca59e44a7877c5dd798a6bdd7d86e
base: origin/main (4d451a8)
captured: 2026-09-25 09:58
check: "For 2026-27, Sundays, 2nd and 4th Saturdays and seeded holidays are non-working days; every other date is a working day"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/auth	2.135s
ok  	officeapp/internal/core/calendar	0.756s
ok  	officeapp/internal/core/clock	0.556s
```

## Story check commands

```
$ go test -count=1 ./internal/core/calendar/   -> exit 0
ok  	officeapp/internal/core/calendar	0.719s
```

## Mutation (scripts/mutate.sh 001)

```
$ scripts/mutate.sh 001   -> PASS
mutation: base origin/main; changed files:
  internal/core/calendar/1790310100_calendar.go
  internal/core/calendar/seed.go
  internal/core/calendar/workdays.go
mutation: killed 6, lived 0, not covered 0
```

## Verdict

check: "For 2026-27, Sundays, 2nd and 4th Saturdays and seeded holidays are non-working days; every other date is a working day"
verdict: PASS — Table tests cover every date in the check (weekly offs, week boundaries, holidays, range, errors); make check green; mutation 6 killed, 0 lived
