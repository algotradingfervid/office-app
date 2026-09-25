# Evidence — story 010

story: factory/stories/010-approvals-approve-reject.md
code-commit: 7b1da8d79a9cacb9b6974c07986d0f4f5e5b2106
base: origin/main (merge-base dfec94a; origin/main now dfec94a)
captured: 2026-09-25 11:01
check: "The current approver can approve as final, which re-runs the form's checks and calls its approval hook exactly once, or reject with a required comment"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	113.691s
ok  	officeapp/internal/core/auth	62.992s
ok  	officeapp/internal/core/calendar	29.779s
ok  	officeapp/internal/core/clock	2.602s
ok  	officeapp/internal/forms/leave	63.580s
```

## Story check commands

```
$ go test -count=1 -run 'ApproveFinal|Reject' ./internal/core/approvals/   -> exit 0
ok  	officeapp/internal/core/approvals	22.343s
```

## Mutation (scripts/mutate.sh 010)

```
$ scripts/mutate.sh 010   -> PASS
mutation: base origin/main; changed files:
  internal/core/approvals/service_decide.go
  internal/core/approvals/service_submit.go
mutation: killed 15, lived 0, not covered 0, timed out 0
```

## Verdict

check: "The current approver can approve as final, which re-runs the form's checks and calls its approval hook exactly once, or reject with a required comment"
verdict: PASS — approve-final re-runs Validate and calls OnFinalApproved once (second approve refused); reject needs a comment; each refusal and rollback tested; 15/15 mutants killed
