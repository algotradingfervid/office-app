# Evidence — story 007

story: factory/stories/007-leave-balance-card.md
code-commit: 722efdae3c211dad9b5792a4cba31105429f261c
base: origin/main (merge-base 1190dde; origin/main now 1190dde)
captured: 2026-09-25 15:45
check: "After signing in, E001 sees cards for CL 4, SL 8 and EL 10.5 with available days on the home page"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	3.224s
ok  	officeapp/internal/core/auth	3.520s
ok  	officeapp/internal/core/calendar	3.193s
ok  	officeapp/internal/core/clock	3.082s
ok  	officeapp/internal/forms/leave	3.888s
```

## Story check commands

```
$ go test -count=1 -run 'BalanceCard' ./internal/forms/leave/   -> exit 0
ok  	officeapp/internal/forms/leave	0.623s
```

## Journey (390 px and 1280 px, light; screen stories only)

```
after (722efda) -> PASS
PASS [390] login E001
PASS [390] see Casual Leave
PASS [390] see 10.5
PASS [390] shot home-balances
PASS [1280] login E001
PASS [1280] see Casual Leave
PASS [1280] see 10.5
PASS [1280] shot home-balances
```
```
before (1190dde), expected to fail where the feature is new:
PASS [390] login E001
FAIL [390] see Casual Leave  -> Locator.wait_for: Timeout 5000ms exceeded.
PASS [1280] login E001
FAIL [1280] see Casual Leave  -> Locator.wait_for: Timeout 5000ms exceeded.
```

Screenshots: `after/` (this branch), `before/` (origin/main).

## Mutation (scripts/mutate.sh 007)

```
$ scripts/mutate.sh 007   -> PASS
mutation: base origin/main; changed files:
  internal/forms/leave/home_card.go
mutation: killed 2, lived 0, not covered 0, timed out 0
```

## Verdict

check: "After signing in, E001 sees cards for CL 4, SL 8 and EL 10.5 with available days on the home page"
verdict: PASS — E001's home shows a Leave balances 2026-27 card: CL 4/4, EL 10.5/10.5, SL 8/8 (balance/available) at 390 and 1280 px; tests pin year and pending subtraction; mutants 2/2 killed
