package leave_test

import (
	"reflect"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	"officeapp/internal/forms/leave"
	"officeapp/internal/testapp"
)

func employeeID(t *testing.T, app *tests.TestApp, code string) string {
	t.Helper()
	emp, err := app.FindFirstRecordByData("users", "employee_code", code)
	if err != nil {
		t.Fatal(err)
	}
	return emp.Id
}

func addEntry(t *testing.T, app *tests.TestApp, empCode, typeCode, year, entryType string, days float64) {
	t.Helper()
	lt, err := app.FindFirstRecordByData("leave_types", "code", typeCode)
	if err != nil {
		t.Fatal(err)
	}
	col, err := app.FindCollectionByNameOrId("leave_ledger")
	if err != nil {
		t.Fatal(err)
	}
	r := core.NewRecord(col)
	r.Set("employee", employeeID(t, app, empCode))
	r.Set("leave_type", lt.Id)
	r.Set("leave_year", year)
	r.Set("entry_type", entryType)
	r.Set("days", days)
	if err := app.Save(r); err != nil {
		t.Fatal(err)
	}
}

func balance(t *testing.T, app *tests.TestApp, empCode, typeCode, year string) float64 {
	t.Helper()
	b, err := leave.Balance(app, employeeID(t, app, empCode), typeCode, year)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestSeededOpeningBalances(t *testing.T) {
	app := testapp.New(t)
	want := map[string]float64{"CL": 4, "SL": 8, "EL": 10.5, "LOP": 0}
	for _, emp := range []string{"E001", "E002", "E003", "E004"} {
		for code, days := range want {
			if got := balance(t, app, emp, code, "2026-27"); got != days {
				t.Errorf("%s %s balance = %v, want %v", emp, code, got, days)
			}
		}
	}
	entries, err := app.FindRecordsByFilter("leave_ledger", "entry_type != 'opening'", "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("seed wrote %d entries that are not opening entries", len(entries))
	}
}

func TestBalanceSumsLedgerEntries(t *testing.T) {
	app := testapp.New(t)
	addEntry(t, app, "E001", "CL", "2026-27", "debit", -2)
	if got := balance(t, app, "E001", "CL", "2026-27"); got != 2 {
		t.Errorf("CL after a -2 debit = %v, want 2", got)
	}

	addEntry(t, app, "E001", "EL", "2026-27", "credit", 1.5)
	addEntry(t, app, "E001", "EL", "2026-27", "debit", -0.5)
	if got := balance(t, app, "E001", "EL", "2026-27"); got != 11.5 {
		t.Errorf("EL 10.5 + 1.5 - 0.5 = %v, want exactly 11.5", got)
	}

	// Other years, types and employees do not count.
	addEntry(t, app, "E001", "CL", "2025-26", "opening", 7)
	addEntry(t, app, "E001", "SL", "2026-27", "debit", -1)
	addEntry(t, app, "E002", "CL", "2026-27", "debit", -3)
	if got := balance(t, app, "E001", "CL", "2026-27"); got != 2 {
		t.Errorf("CL 2026-27 = %v, want 2", got)
	}
	if got := balance(t, app, "E001", "CL", "2025-26"); got != 7 {
		t.Errorf("CL 2025-26 = %v, want 7", got)
	}
	if got := balance(t, app, "E001", "CL", "2027-28"); got != 0 {
		t.Errorf("CL 2027-28 with no entries = %v, want 0", got)
	}
}

func TestBalancesListsActiveQuotaTypes(t *testing.T) {
	app := testapp.New(t)
	emp := employeeID(t, app, "E001")
	addEntry(t, app, "E001", "SL", "2026-27", "debit", -1)

	got, err := leave.Balances(app, emp, "2026-27")
	if err != nil {
		t.Fatal(err)
	}
	want := []leave.TypeBalance{
		{Code: "CL", Name: "Casual Leave", Balance: 4},
		{Code: "EL", Name: "Earned Leave", Balance: 10.5},
		{Code: "SL", Name: "Sick Leave", Balance: 7},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Balances = %+v, want %+v", got, want)
	}

	cl, err := app.FindFirstRecordByData("leave_types", "code", "CL")
	if err != nil {
		t.Fatal(err)
	}
	cl.Set("active", false)
	if err := app.Save(cl); err != nil {
		t.Fatal(err)
	}
	got, err = leave.Balances(app, emp, "2026-27")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want[1:]) {
		t.Errorf("Balances with CL inactive = %+v, want %+v", got, want[1:])
	}
}
