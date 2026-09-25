# Evidence — story 005

story: factory/stories/005-leave-balances.md
code-commit: 8d5a07a3321c9bbe89b1ad39098a7402859dbcc1
base: origin/main (6cebbb9)
captured: 2026-09-25 10:24
check: "An employee's balance is the sum of their ledger entries for a type and leave year"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	2.027s
ok  	officeapp/internal/core/auth	6.685s
ok  	officeapp/internal/core/calendar	4.085s
ok  	officeapp/internal/core/clock	2.210s
ok  	officeapp/internal/forms/leave	4.689s
```

## Story check commands

```
$ go test -count=1 -run 'Balance' ./internal/forms/leave/   -> exit 0
ok  	officeapp/internal/forms/leave	2.261s
```

## Mutation (scripts/mutate.sh 005)

```
$ scripts/mutate.sh 005   -> PASS
mutation: base origin/main; changed files:
  internal/forms/leave/balance.go
  internal/forms/leave/seed_balances.go
mutation: killed 7, lived 0, not covered 0, timed out 0
```

## Verdict

check: "An employee's balance is the sum of their ledger entries for a type and leave year"
verdict: PASS — Balance sums ledger entries per employee/type/year (exact halves, other years/types/employees excluded); Balances lists active quota types; demo openings seeded; 7/7 mutants killed
