package leave_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	"officeapp/internal/forms/leave"
	"officeapp/internal/testapp"
)

func TestLeaveYear(t *testing.T) {
	cases := map[string]string{
		"2026-10-12": "2026-27",
		"2027-03-31": "2026-27",
		"2027-04-01": "2027-28",
		"2026-04-01": "2026-27",
		"2026-03-31": "2025-26",
		"2099-12-31": "2099-00",
	}
	for date, want := range cases {
		if got := leave.LeaveYear(date); got != want {
			t.Errorf("LeaveYear(%s) = %s, want %s", date, got, want)
		}
	}
}

func TestCollectionsClosedToAPI(t *testing.T) {
	app := testapp.New(t)
	for _, name := range []string{"leave_types", "leave_rules", "leave_requests", "leave_ledger"} {
		c, err := app.FindCollectionByNameOrId(name)
		if err != nil {
			t.Fatalf("collection %s: %v", name, err)
		}
		if c.ListRule != nil || c.ViewRule != nil || c.CreateRule != nil || c.UpdateRule != nil || c.DeleteRule != nil {
			t.Errorf("%s API rules must all be nil", name)
		}
	}
}

func TestSeededTypesAndRules(t *testing.T) {
	app := testapp.New(t)
	types, err := app.FindAllRecords("leave_types")
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 9 {
		t.Fatalf("got %d leave types, want 9", len(types))
	}
	for _, lt := range types {
		want := lt.GetString("code") != "CO"
		if lt.GetBool("active") != want {
			t.Errorf("%s active = %v, want %v", lt.GetString("code"), lt.GetBool("active"), want)
		}
	}
	type want struct {
		credit, count                       string
		perYear, monthly, maxRun, yearlyCap float64
		maxTimes, notice, backdate          int
		halfDay, encash, probation          bool
		attachAfter, carryCap               float64
	}
	wants := map[string]want{
		"CL":  {"monthly", "working_days", 12, 1, 2, 0, 0, 1, 0, true, false, true, 0, 0},
		"SL":  {"yearly", "working_days", 8, 0, 0, 0, 0, 0, 7, true, false, true, 2, 0},
		"EL":  {"monthly", "working_days", 18, 1.5, 0, 0, 0, 7, 0, true, true, false, 0, 30},
		"ML":  {"per_event", "calendar_days", 0, 0, 182, 0, 0, 0, 0, false, false, true, 0, 0},
		"PL":  {"per_event", "working_days", 0, 0, 5, 0, 0, 0, 0, false, false, true, 0, 0},
		"BL":  {"per_event", "working_days", 0, 0, 3, 0, 0, 0, 0, false, false, true, 0, 0},
		"MRL": {"per_event", "working_days", 0, 0, 3, 0, 1, 0, 0, false, false, true, 0, 0},
		"CO":  {"per_event", "working_days", 0, 0, 0, 0, 0, 0, 0, false, false, true, 0, 0},
		"LOP": {"per_event", "working_days", 0, 0, 5, 15, 0, 0, 0, false, false, true, 0, 0},
	}
	for code, w := range wants {
		r, err := leave.RuleFor(app, code, "2026-10-12")
		if err != nil {
			t.Fatalf("RuleFor(%s): %v", code, err)
		}
		got := want{
			r.GetString("credit_method"), r.GetString("count_mode"),
			r.GetFloat("days_per_year"), r.GetFloat("monthly_credit"), r.GetFloat("max_consecutive_days"), r.GetFloat("yearly_cap"),
			r.GetInt("max_times_per_employment"), r.GetInt("min_notice_days"), r.GetInt("max_backdate_days"),
			r.GetBool("half_day_allowed"), r.GetBool("encashable"), r.GetBool("allowed_in_probation"),
			r.GetFloat("attachment_after_days"), r.GetFloat("carry_forward_cap"),
		}
		if got != w {
			t.Errorf("%s rule = %+v, want %+v", code, got, w)
		}
		if r.GetString("effective_from") != "2026-04-01" {
			t.Errorf("%s effective_from = %s", code, r.GetString("effective_from"))
		}
	}
	if _, err := leave.RuleFor(app, "CL", "2026-03-31"); err == nil {
		t.Error("RuleFor before the first rule found one")
	}
	if _, err := leave.RuleFor(app, "XX", "2026-10-12"); err == nil {
		t.Error("RuleFor an unknown type found one")
	}
}

func TestRuleForPicksVersion(t *testing.T) {
	app := testapp.New(t)
	first, err := leave.RuleFor(app, "CL", "2026-10-12")
	if err != nil {
		t.Fatal(err)
	}
	second := core.NewRecord(first.Collection())
	second.Load(first.FieldsData())
	second.Id = ""
	second.Set("effective_from", "2026-12-01")
	second.Set("max_consecutive_days", 3)
	if err := app.Save(second); err != nil {
		t.Fatal(err)
	}
	cases := map[string]string{
		"2026-04-01": first.Id,
		"2026-11-30": first.Id,
		"2026-12-01": second.Id,
		"2027-05-01": second.Id,
	}
	for date, want := range cases {
		r, err := leave.RuleFor(app, "CL", date)
		if err != nil {
			t.Fatalf("RuleFor(CL, %s): %v", date, err)
		}
		if r.Id != want {
			t.Errorf("RuleFor(CL, %s) = %s (from %s), want %s", date, r.Id, r.GetString("effective_from"), want)
		}
	}
}

func ledgerEntry(t *testing.T, app *tests.TestApp, entryType, periodKey string, days float64) *core.Record {
	t.Helper()
	emp, err := app.FindFirstRecordByData("users", "employee_code", "E001")
	if err != nil {
		t.Fatal(err)
	}
	cl, err := app.FindFirstRecordByData("leave_types", "code", "CL")
	if err != nil {
		t.Fatal(err)
	}
	col, err := app.FindCollectionByNameOrId("leave_ledger")
	if err != nil {
		t.Fatal(err)
	}
	r := core.NewRecord(col)
	r.Set("employee", emp.Id)
	r.Set("leave_type", cl.Id)
	r.Set("leave_year", "2026-27")
	r.Set("entry_type", entryType)
	r.Set("days", days)
	r.Set("period_key", periodKey)
	return r
}

func TestLedgerAppendOnly(t *testing.T) {
	app := testapp.New(t)
	e := ledgerEntry(t, app, "credit", "2026-11", 1)
	if err := app.Save(e); err != nil {
		t.Fatal(err)
	}
	e.Set("days", 5)
	if err := app.Save(e); err == nil {
		t.Error("a ledger entry was updated")
	}
	if err := app.Delete(e); err == nil {
		t.Error("a ledger entry was deleted")
	}
	stored, err := app.FindRecordById("leave_ledger", e.Id)
	if err != nil {
		t.Fatal(err)
	}
	if stored.GetFloat("days") != 1 {
		t.Errorf("stored days = %v, want 1", stored.GetFloat("days"))
	}
}

func TestLedgerPeriodKeyUnique(t *testing.T) {
	app := testapp.New(t)
	if err := app.Save(ledgerEntry(t, app, "credit", "2026-11", 1)); err != nil {
		t.Fatal(err)
	}
	if err := app.Save(ledgerEntry(t, app, "credit", "2026-11", 1)); err == nil {
		t.Error("a second credit for 2026-11 was saved")
	}
	if err := app.Save(ledgerEntry(t, app, "credit", "2026-12", 1)); err != nil {
		t.Errorf("credit for another period: %v", err)
	}
	for i := 0; i < 2; i++ {
		if err := app.Save(ledgerEntry(t, app, "adjustment", "", 0.5)); err != nil {
			t.Errorf("adjustment %d without period key: %v", i+1, err)
		}
	}
}

func TestDaysAreHalves(t *testing.T) {
	app := testapp.New(t)
	cases := map[float64]bool{0.5: true, -1.5: true, 2: true, 0.25: false, 1.3: false, -0.75: false}
	for days, ok := range cases {
		err := app.Save(ledgerEntry(t, app, "adjustment", "", days))
		if (err == nil) != ok {
			t.Errorf("ledger days %v: err = %v, want ok = %v", days, err, ok)
		}
	}
	rule, err := leave.RuleFor(app, "EL", "2026-10-12")
	if err != nil {
		t.Fatal(err)
	}
	rule.Set("monthly_credit", 1.25)
	if err := app.Save(rule); err == nil {
		t.Error("a rule with monthly_credit 1.25 was saved")
	}

	col, err := app.FindCollectionByNameOrId("leave_requests")
	if err != nil {
		t.Fatal(err)
	}
	for days, ok := range map[float64]bool{1.5: true, 1.25: false} {
		req := core.NewRecord(col)
		req.Set("leave_type", rule.GetString("leave_type"))
		req.Set("from_date", "2026-10-12")
		req.Set("from_session", "full")
		req.Set("to_date", "2026-10-13")
		req.Set("to_session", "first_half")
		req.Set("days", days)
		req.Set("leave_year", "2026-27")
		req.Set("rule", rule.Id)
		if err := app.Save(req); (err == nil) != ok {
			t.Errorf("leave request days %v: err = %v, want ok = %v", days, err, ok)
		}
	}
}
