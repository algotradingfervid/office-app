# Evidence — story 003

story: factory/stories/003-leave-contract.md
code-commit: 09509d8d2642886e62f3517edee57380c6c6697f
base: origin/main (4d451a8)
captured: 2026-09-25 10:01
check: "The leave collections exist with the design's fields, the default rules for all nine leave types are seeded, and a ledger entry can never be edited or deleted"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/auth	1.916s
ok  	officeapp/internal/core/clock	0.860s
ok  	officeapp/internal/forms/leave	0.728s
```

## Story check commands

```
$ go test -count=1 ./internal/forms/leave/   -> exit 0
ok  	officeapp/internal/forms/leave	0.693s
```

## Mutation (scripts/mutate.sh 003)

```
$ scripts/mutate.sh 003   -> PASS
mutation: base origin/main; changed files:
  internal/forms/leave/1790310300_leave.go
  internal/forms/leave/ledger_guard.go
  internal/forms/leave/rules.go
  internal/forms/leave/seed_rules.go
mutation: killed 11, lived 0, not covered 0
```

## Verdict

check: "The leave collections exist with the design's fields, the default rules for all nine leave types are seeded, and a ledger entry can never be edited or deleted"
verdict: PASS — make check green; story tests cover LeaveYear bounds, RuleFor versions, seeded 9 types/rules, ledger append-only, period_key partial unique index, half-day fields; mutation 11 killed, 0 lived. request relations deferred to story 021 by floor decision
