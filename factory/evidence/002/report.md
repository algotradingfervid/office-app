# Evidence — story 002

story: factory/stories/002-approvals-contract.md
code-commit: 85fff96ac8929d61ba3f5298b04ff11019e99b96
base: origin/main (4d451a8)
captured: 2026-09-25 09:59
check: "A request's status can only move along the design's transition table, and an approval step can never be edited or deleted"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	4.041s
ok  	officeapp/internal/core/auth	6.697s
ok  	officeapp/internal/core/clock	0.856s
```

## Story check commands

```
$ go test -count=1 ./internal/core/approvals/   -> exit 0
ok  	officeapp/internal/core/approvals	3.334s
```

## Mutation (scripts/mutate.sh 002)

```
$ scripts/mutate.sh 002   -> PASS
mutation: base origin/main; changed files:
  internal/core/approvals/1790310200_requests.go
  internal/core/approvals/forms.go
  internal/core/approvals/history_guard.go
  internal/core/approvals/transitions.go
mutation: killed 3, lived 0, not covered 0
```

## Verdict

check: "A request's status can only move along the design's transition table, and an approval step can never be edited or deleted"
verdict: PASS — Every §7 row allowed and all 54 other status/action pairs refused; steps save, update and delete fail; RegisterForm/FormFor round-trip; mutation 3 killed 0 lived
