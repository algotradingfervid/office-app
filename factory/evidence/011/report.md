# Evidence — story 011

story: factory/stories/011-approvals-my-requests-page.md
code-commit: 6166794f412c227f0339b167ebf098238c4aa4c8
base: origin/main (merge-base dfec94a; origin/main now dfec94a)
captured: 2026-09-25 11:04
check: "A signed-in employee sees a list of only their own requests with form summary, dates submitted and status, newest first"
verdict: PASS

## make check

```
$ make check   -> PASS
0 issues.
imports: ok
ok  	officeapp/internal/core/approvals	26.126s
ok  	officeapp/internal/core/auth	29.103s
ok  	officeapp/internal/core/calendar	10.989s
ok  	officeapp/internal/core/clock	1.705s
ok  	officeapp/internal/forms/leave	14.864s
```

## Story check commands

```
$ go test -count=1 -run 'MyRequests' ./internal/core/approvals/   -> exit 0
ok  	officeapp/internal/core/approvals	27.177s
```

## Journey (390 px and 1280 px, light; screen stories only)

```
after (d250af6) -> PASS
PASS [390] login E001
PASS [390] goto /requests
PASS [390] see My requests
PASS [390] see You have not made any requests yet.
PASS [390] shot my-requests
PASS [1280] login E001
PASS [1280] goto /requests
PASS [1280] see My requests
PASS [1280] see You have not made any requests yet.
PASS [1280] shot my-requests
```
```
before (dfec94a), expected to fail where the feature is new:
PASS [390] login E001
PASS [390] goto /requests
FAIL [390] see My requests  -> Locator.wait_for: Timeout 5000ms exceeded.
FAIL [390] console  -> Failed to load resource: the server responded with a status of 404 (Not Found)
PASS [1280] login E001
PASS [1280] goto /requests
FAIL [1280] see My requests  -> Locator.wait_for: Timeout 5000ms exceeded.
FAIL [1280] console  -> Failed to load resource: the server responded with a status of 404 (Not Found)
```

Screenshots: `after/` (this branch), `before/` (origin/main).

## Mutation (scripts/mutate.sh 011)

```
$ scripts/mutate.sh 011   -> PASS
mutation: base origin/main; changed files:
  internal/core/approvals/my_requests_page.go
mutation: killed 6, lived 0, not covered 0, timed out 0
```

## Verdict

check: "A signed-in employee sees a list of only their own requests with form summary, dates submitted and status, newest first"
verdict: PASS — GET /requests lists only the signed-in employee's requests newest first with summary, IST date, status and approver (ApiScenario tests); journey shows the empty state at 390/1280; mutation 6 killed 0 lived
