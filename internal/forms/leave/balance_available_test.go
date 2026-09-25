package leave_test

import (
	"reflect"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	"officeapp/internal/core/approvals"
	"officeapp/internal/forms/leave"
	"officeapp/internal/testapp"
)

// requestLeave saves a request with the given status and its leave request of typeCode, leave year and days.
func requestLeave(t *testing.T, app *tests.TestApp, empCode string, status approvals.Status, typeCode, year string, days float64) {
	t.Helper()
	col, err := app.FindCollectionByNameOrId("requests")
	if err != nil {
		t.Fatal(err)
	}
	req := core.NewRecord(col)
	req.Set("form_type", "leave")
	req.Set("requester", employeeID(t, app, empCode))
	req.Set("status", string(status))
	req.Set("submitted_at", "2026-10-01 04:30:00.000Z")
	req.Set("form_version", 1)
	if err := app.Save(req); err != nil {
		t.Fatal(err)
	}
	rule, err := leave.RuleFor(app, typeCode, "2026-10-12")
	if err != nil {
		t.Fatal(err)
	}
	col, err = app.FindCollectionByNameOrId("leave_requests")
	if err != nil {
		t.Fatal(err)
	}
	lr := core.NewRecord(col)
	lr.Set("request", req.Id)
	lr.Set("leave_type", rule.GetString("leave_type"))
	lr.Set("from_date", "2026-10-12")
	lr.Set("from_session", "full")
	lr.Set("to_date", "2026-10-12")
	lr.Set("to_session", "full")
	lr.Set("days", days)
	lr.Set("leave_year", year)
	lr.Set("rule", rule.Id)
	if err := app.Save(lr); err != nil {
		t.Fatal(err)
	}
}

func available(t *testing.T, app *tests.TestApp, empCode, typeCode, year string) float64 {
	t.Helper()
	a, err := leave.Available(app, employeeID(t, app, empCode), typeCode, year)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func TestAvailableSubtractsPendingDays(t *testing.T) {
	app := testapp.New(t)
	if got := available(t, app, "E001", "CL", "2026-27"); got != 4 {
		t.Errorf("CL available with nothing pending = %v, want 4", got)
	}
	requestLeave(t, app, "E001", approvals.Pending, "CL", "2026-27", 1)
	if got := available(t, app, "E001", "CL", "2026-27"); got != 3 {
		t.Errorf("CL available with a pending 1-day CL = %v, want 3", got)
	}
	if got := balance(t, app, "E001", "CL", "2026-27"); got != 4 {
		t.Errorf("CL balance with a pending 1-day CL = %v, want 4", got)
	}
	requestLeave(t, app, "E001", approvals.Pending, "CL", "2026-27", 0.5)
	if got := available(t, app, "E001", "CL", "2026-27"); got != 2.5 {
		t.Errorf("CL available with 1 + 0.5 days pending = %v, want 2.5", got)
	}
}

func TestAvailableIgnoresClosedRequests(t *testing.T) {
	app := testapp.New(t)
	for _, s := range []approvals.Status{approvals.Approved, approvals.Rejected, approvals.Cancelled, approvals.CancelRequested} {
		requestLeave(t, app, "E001", s, "CL", "2026-27", 1)
	}
	if got := available(t, app, "E001", "CL", "2026-27"); got != 4 {
		t.Errorf("CL available with only closed requests = %v, want 4", got)
	}
}

func TestAvailableCountsOnlyThatTypeYearAndEmployee(t *testing.T) {
	app := testapp.New(t)
	requestLeave(t, app, "E001", approvals.Pending, "EL", "2026-27", 2)
	requestLeave(t, app, "E001", approvals.Pending, "CL", "2025-26", 1)
	requestLeave(t, app, "E002", approvals.Pending, "CL", "2026-27", 1)
	if got := available(t, app, "E001", "CL", "2026-27"); got != 4 {
		t.Errorf("E001 CL 2026-27 available = %v, want 4", got)
	}
	if got := available(t, app, "E001", "EL", "2026-27"); got != 8.5 {
		t.Errorf("E001 EL available with 2 days pending = %v, want 8.5", got)
	}
}

func TestAvailableBalancesListsBalanceAndAvailable(t *testing.T) {
	app := testapp.New(t)
	requestLeave(t, app, "E001", approvals.Pending, "CL", "2026-27", 1)
	got, err := leave.AvailableBalances(app, employeeID(t, app, "E001"), "2026-27")
	if err != nil {
		t.Fatal(err)
	}
	want := []leave.TypeAvailable{
		{TypeBalance: leave.TypeBalance{Code: "CL", Name: "Casual Leave", Balance: 4}, Available: 3},
		{TypeBalance: leave.TypeBalance{Code: "EL", Name: "Earned Leave", Balance: 10.5}, Available: 10.5},
		{TypeBalance: leave.TypeBalance{Code: "SL", Name: "Sick Leave", Balance: 8}, Available: 8},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AvailableBalances = %+v, want %+v", got, want)
	}
}
