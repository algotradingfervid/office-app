# Evidence — story 009

story: factory/stories/009-leave-limit-checks.md
code-commit: dbd086e5d620c52d9445ec1eb0eb37a3f3bac0ff
base: origin/main (merge-base 1190dde; origin/main now 1190dde)
captured: 2026-09-25 15:53
check: "A leave request is refused when it exceeds the consecutive-day limit (adjacent requests included), the available balance, the yearly cap, or the once-per-employment limit, or lacks a required attachment; LOP while paid leave remains is a warning"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	1.247s
ok  	officeapp/internal/core/auth	1.623s
ok  	officeapp/internal/core/calendar	1.317s
ok  	officeapp/internal/core/clock	2.096s
ok  	officeapp/internal/forms/leave	3.159s
```

## Story check commands

```
$ go test -count=1 -run 'LimitChecks' ./internal/forms/leave/   -> exit 0
ok  	officeapp/internal/forms/leave	0.881s
```

## Mutation (scripts/mutate.sh 009)

```
$ scripts/mutate.sh 009   -> PASS
mutation: base origin/main; changed files:
  internal/forms/leave/policy_limits.go
mutation: killed 60, lived 0, not covered 0, timed out 0
```

## Verdict

check: "A leave request is refused when it exceeds the consecutive-day limit (adjacent requests included), the available balance, the yearly cap, or the once-per-employment limit, or lacks a required attachment; LOP while paid leave remains is a warning"
verdict: PASS — Every check-6-to-11 case in the story holds (50+ table cases, incl. adjacency, balance, LOP cap 15, MRL once, SL attachment, LOP warning); make check green; mutation 60 killed, 0 lived
