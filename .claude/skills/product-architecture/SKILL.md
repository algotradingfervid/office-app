---
name: product-architecture
description: "Office App code shape: where each kind of Go code goes (module, layer, file), how a feature registers itself, how migrations, templates, transactions and tests are written. Load BEFORE writing, changing or reviewing any application code in this repo (Go, templates, migrations, tests), even for a one-line change."
user-invocable: false
---

# Office App architecture — how to write code here

Full reasoning: `factory/ARCHITECTURE.md`. Business rules: design spec `docs/superpowers/specs/2026-09-25-office-app-leave-design.md`. Words: `factory/GLOSSARY.md`.

## Where code goes

| You are writing | It goes in |
|---|---|
| A shared platform feature (users, approvals, calendar, notifications, audit) | `internal/core/<module>/` |
| Anything specific to one form | `internal/forms/<form>/` (leave: `internal/forms/leave/`) |
| A collection or field | a new migration file in the owning module: `<unix-time>_<what>.go` |
| A page | handler in `<module>/<feature>_page.go`, templates in `<module>/templates/<page>*.html`, embedded by that page file only |
| An htmx fragment | its own template file `<module>/templates/<page>_<fragment>.html` |
| Business rules (checks, counting, balances) | plain functions/services in the module; no HTTP types in them |
| Demo data | the module's seeder, registered with `seed.Add` |
| Helpers another module needs | exported from the module's **root package** only |
| Code only this module uses | unexported, or in `<module>/internal/...` |

Never edit `internal/modules/modules.go` or a module's root `Register` function to add a feature. Only a story that creates a whole new module touches `modules.go` (one line).

## Layers inside a module

`handler (routes, form parsing, templates) → service (rules, transaction) → data (PocketBase records)`

- Handlers parse input, call **one** service function, render. They never `Save` records.
- Services own the rules and the transaction. Signature style: `func Submit(app core.App, c clock.Clock, in SubmitInput) (Result, error)`; inside, `app.RunInTransaction(func(txApp core.App) error { ... })` and **every** read and write in it uses `txApp`.
- Data access is PocketBase's API (`FindFirstRecordByData`, `FindRecordsByFilter` with `dbx.Params`, `Save`). No raw SQL strings built from input.
- Time comes from a `clock.Clock` argument (`clock.System{}` in handlers, `clock.Fixed(...)` in tests). Dates are `YYYY-MM-DD` strings; compare them as strings or parse with `time.ParseInLocation(time.DateOnly, s, clock.IST)`.

## How a feature registers itself (the pattern every file copies)

```go
// internal/core/calendar/holidays_page.go
package calendar

func init() {
	addPart(func(app core.App) {
		app.OnServe().BindFunc(func(se *core.ServeEvent) error {
			se.Router.GET("/holidays", showHolidays).BindFunc(auth.RequireUser)
			return se.Next()
		})
		seed.Add(app, "calendar.holidays", seedHolidays) // only if this file seeds data
	})
}
```

Home cards: `home.AddCard(app, func(e *core.RequestEvent) (template.HTML, error) { ... })` from the form's own file.

## Worked example: one feature across the layers

Feature: "the calendar knows holidays; the apply page asks whether a date is a working day".

1. **Migration** — `internal/core/calendar/1790301000_holidays.go`:
```go
func init() {
	m.Register(func(app core.App) error {
		c := core.NewBaseCollection("holidays")
		// ListRule..DeleteRule stay nil: no REST access
		c.Fields.Add(
			&core.TextField{Name: "date", Required: true, Pattern: `^\d{4}-\d{2}-\d{2}$`},
			&core.TextField{Name: "name", Required: true, Max: 100},
		)
		c.AddIndex("idx_holidays_date", true, "date", "")
		return app.Save(c)
	}, nil)
}
```
2. **Service** — `internal/core/calendar/workdays.go`: `func IsWorkingDay(app core.App, date string) (bool, error)` reads the weekly-off setting and `holidays` through the `app` it is given (a `txApp` when called inside a transaction).
3. **Handler** — a page file with `init(){ addPart(...) }`, parses the form, calls the service, renders. Each page file embeds **only its own templates**, so two page stories never share a file:
```go
//go:embed templates/holidays.html templates/holidays_row.html
var holidaysFS embed.FS
var holidaysPages = web.NewPages(holidaysFS)
// ... holidaysPages.Render(e, "holidays", "Holidays", data)
```
   A non-page fragment (e.g. a home card) parses its own template with `template.ParseFS` and returns `template.HTML`.
4. **Template** — `internal/core/calendar/templates/holidays.html` defines `title` and `content`; uses `.Data`, `.User`; no inline `<script>`/`style=`.
5. **Tests** — `workdays_test.go` (table-driven, `clock.Fixed`) for the rule; a page test with `tests.ApiScenario{..., TestAppFactory: testapp.Factory}` and `Headers: testapp.FormHeaders(testapp.SessionCookie(t, app, "E004"))` for POSTs.

## Rules for tests

- Business rules: table-driven unit tests with explicit dates (e.g. `2026-10-10` is a 2nd Saturday → weekly off).
- Anything touching records: `app := testapp.New(t)` — migrated, demo data seeded (E001–E004), own DB copy.
- Pages: `tests.ApiScenario` with `TestAppFactory: testapp.Factory`, or a closure returning a prepared app.
- A test must fail if the feature is removed. Mutation testing (`scripts/mutate.sh`) checks this.

## Templates and UI

- Stock Pico styling until the design-system retrofit; semantic HTML (`article`, `form`, `label`, `table`), no custom visual work.
- htmx: `hx-post`/`hx-get` + `hx-target`/`hx-swap`; state changes are POST; responses are server-rendered fragments. No `hx-on`, no inline JS (CSP).
- User-facing messages are plain sentences ("CL allows at most 2 consecutive days."), never Go error text.

## Anti-patterns (reviewers block these)

- Editing `modules.go` / `Register` instead of `addPart`.
- Outer `app` inside `RunInTransaction`; two transactions for one user action.
- `time.Now()` in rules; PocketBase `date` fields for calendar dates.
- Collection API rules other than `nil`; calling `/api/collections/...` from pages.
- Updating/deleting `approval_steps`, `leave_ledger`, `audit_log` rows.
- Editing a migration already on `main` (a hook blocks it) — add a new one.
- Importing another module's subpackage, or core importing a form (`make check` fails).
- New dependencies for what a short tested function does; abstractions with one caller.
