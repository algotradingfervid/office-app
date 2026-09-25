# Evidence — story 006

story: factory/stories/006-approvals-submit-forward.md
code-commit: 4fc26bbe638f68382942debf21fd90535fa69c1e
base: origin/main (6cebbb9)
captured: 2026-09-25 10:25
check: "Submitting creates a pending request with a submitted step for the chosen approver; the current approver can forward it to another active approver, at most 5 times"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	6.366s
ok  	officeapp/internal/core/auth	5.131s
ok  	officeapp/internal/core/calendar	2.637s
ok  	officeapp/internal/core/clock	0.911s
ok  	officeapp/internal/forms/leave	3.938s
```

## Story check commands

```
$ go test -count=1 -run 'Submit|Forward' ./internal/core/approvals/   -> exit 0
ok  	officeapp/internal/core/approvals	8.427s
```

## Mutation (scripts/mutate.sh 006)

```
$ scripts/mutate.sh 006   -> PASS
mutation: base origin/main; changed files:
  internal/core/approvals/service_submit.go
mutation: killed 23, lived 0, not covered 0, timed out 0
```

## Verdict

check: "Submitting creates a pending request with a submitted step for the chosen approver; the current approver can forward it to another active approver, at most 5 times"
verdict: PASS — Submit/forward tests pass for every listed case incl. per-role actor checks; make check green; 23/23 mutants killed
