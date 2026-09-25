package leave_test

import (
	"errors"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/filesystem"

	"officeapp/internal/core/approvals"
	"officeapp/internal/forms/leave"
	"officeapp/internal/testapp"
)

// October 2026: 5 Mon … 9 Fri, 10 (2nd Sat) and 11 (Sun) off, 12 Mon … 17 Sat (working), 24 (4th Sat) off.
// 2 October (Fri) is a holiday; 3 October is the 1st Saturday, a working day.
type limitRow struct {
	emp    string
	status approvals.Status
	code   string
	from   string
	fs     string
	to     string
	ts     string
	days   float64
	attach bool
}

// limitLeave saves a request with its leave request, as SubmitRequest does before it calls Validate.
func limitLeave(t *testing.T, app *tests.TestApp, l limitRow) *core.Record {
	t.Helper()
	rule, err := leave.RuleFor(app, l.code, "2026-10-12") // the seeded rules start on 2026-04-01
	if err != nil {
		t.Fatal(err)
	}
	requests, err := app.FindCollectionByNameOrId("requests")
	if err != nil {
		t.Fatal(err)
	}
	req := core.NewRecord(requests)
	req.Set("form_type", "leave")
	req.Set("requester", employeeID(t, app, l.emp))
	req.Set("status", string(l.status))
	req.Set("submitted_at", "2026-10-01 04:30:00.000Z")
	req.Set("form_version", 1)
	if err := app.Save(req); err != nil {
		t.Fatal(err)
	}
	col, err := app.FindCollectionByNameOrId("leave_requests")
	if err != nil {
		t.Fatal(err)
	}
	lr := core.NewRecord(col)
	lr.Set("request", req.Id)
	lr.Set("leave_type", rule.GetString("leave_type"))
	lr.Set("from_date", l.from)
	lr.Set("from_session", l.fs)
	lr.Set("to_date", l.to)
	lr.Set("to_session", l.ts)
	lr.Set("days", l.days)
	lr.Set("leave_year", leave.LeaveYear(l.from))
	lr.Set("rule", rule.Id)
	if l.attach {
		f, err := filesystem.NewFileFromBytes([]byte("%PDF-1.4\n%%EOF\n"), "certificate.pdf")
		if err != nil {
			t.Fatal(err)
		}
		lr.Set("attachment", f)
	}
	if err := app.Save(lr); err != nil {
		t.Fatal(err)
	}
	return lr
}

// mine is an E001 request of full days with the given status.
func mine(status approvals.Status, code, from, to string, days float64) limitRow {
	return limitRow{"E001", status, code, from, "full", to, "full", days, false}
}

func req(code, from, to string, days float64) limitRow {
	return mine(approvals.Pending, code, from, to, days)
}

func TestLimitChecks(t *testing.T) {
	const (
		pend     = approvals.Pending
		appr     = approvals.Approved
		consec   = "Casual Leave allows at most 2 consecutive days."
		lopWarn  = "Loss of Pay applied for while paid leave is available (Casual Leave: 4, Earned Leave: 10.5)."
		lopCap   = "Loss of Pay is limited to 15 days in a leave year. Contact HR."
		pl       = "Paternity Leave allows at most 5 consecutive days."
		sickNote = "Sick Leave of more than 2 days needs an attachment."
		mrlOnce  = "Marriage Leave can be taken only once during employment."
	)
	sl3 := req("SL", "2026-10-12", "2026-10-14", 3)
	sl3.attach = true
	type adjustment struct {
		code string
		days float64
	}
	cases := []struct {
		name     string
		existing []limitRow
		ledger   []adjustment   // E001 2026-27 adjustment entries
		rule     map[string]any // changes to the request type's rule
		req      limitRow
		message  string // "" = allowed
		warning  string
	}{
		{name: "inactive type refused", req: req("CO", "2026-10-12", "2026-10-12", 1),
			message: "Compensatory Off cannot be applied for at present."},
		{name: "days come from the day count", req: req("CL", "2026-10-10", "2026-10-12", 1),
			message: "Leave must start and end on a working day."},

		// 6. consecutive days, adjacent requests of the same type added together
		{name: "CL 2 days allowed", req: req("CL", "2026-10-12", "2026-10-13", 2)},
		{name: "CL 3 days refused", req: req("CL", "2026-10-12", "2026-10-14", 3), message: consec},
		{name: "CL Mon after 2 CL Thu–Fri refused", existing: []limitRow{req("CL", "2026-10-08", "2026-10-09", 2)},
			req: req("CL", "2026-10-12", "2026-10-12", 1), message: consec},
		{name: "CL Mon after 1 CL Fri allowed", existing: []limitRow{req("CL", "2026-10-09", "2026-10-09", 1)},
			req: req("CL", "2026-10-12", "2026-10-12", 1)},
		{name: "CL Fri before 2 CL Mon–Tue refused", existing: []limitRow{req("CL", "2026-10-12", "2026-10-13", 2)},
			req: req("CL", "2026-10-09", "2026-10-09", 1), message: consec},
		{name: "CL Sat after CL Wed–Thu across a holiday refused", existing: []limitRow{req("CL", "2026-09-30", "2026-10-01", 2)},
			req: req("CL", "2026-10-03", "2026-10-03", 1), message: consec},
		{name: "working day between is not adjacent (before)", existing: []limitRow{req("CL", "2026-10-07", "2026-10-08", 2)},
			req: req("CL", "2026-10-12", "2026-10-12", 1)},
		{name: "working day between is not adjacent (after)", existing: []limitRow{req("CL", "2026-10-12", "2026-10-13", 2)},
			req: req("CL", "2026-10-08", "2026-10-08", 1)},
		{name: "approved adjacent CL counts", existing: []limitRow{mine(appr, "CL", "2026-10-08", "2026-10-09", 2)},
			req: req("CL", "2026-10-12", "2026-10-12", 1), message: consec},
		{name: "cancel_requested adjacent CL counts", existing: []limitRow{mine(approvals.CancelRequested, "CL", "2026-10-08", "2026-10-09", 2)},
			req: req("CL", "2026-10-12", "2026-10-12", 1), message: consec},
		{name: "rejected adjacent CL does not count", existing: []limitRow{mine(approvals.Rejected, "CL", "2026-10-08", "2026-10-09", 2)},
			req: req("CL", "2026-10-12", "2026-10-12", 1)},
		{name: "cancelled adjacent CL does not count", existing: []limitRow{mine(approvals.Cancelled, "CL", "2026-10-08", "2026-10-09", 2)},
			req: req("CL", "2026-10-12", "2026-10-12", 1)},
		{name: "adjacent SL does not count for CL", existing: []limitRow{req("SL", "2026-10-08", "2026-10-09", 2)},
			req: req("CL", "2026-10-12", "2026-10-12", 1)},
		{name: "another employee's adjacent CL does not count",
			existing: []limitRow{{"E002", pend, "CL", "2026-10-08", "full", "2026-10-09", "full", 2, false}},
			req:      req("CL", "2026-10-12", "2026-10-12", 1)},
		{name: "worked afternoon before breaks adjacency",
			existing: []limitRow{{"E001", pend, "CL", "2026-10-08", "full", "2026-10-09", "first_half", 1.5, false}},
			req:      req("CL", "2026-10-12", "2026-10-12", 1)},
		{name: "worked morning after breaks adjacency", existing: []limitRow{req("CL", "2026-10-08", "2026-10-09", 2)},
			req: limitRow{"E001", pend, "CL", "2026-10-12", "second_half", "2026-10-12", "second_half", 0.5, false}},
		{name: "worked morning of the next request breaks adjacency",
			existing: []limitRow{{"E001", pend, "CL", "2026-10-12", "second_half", "2026-10-13", "full", 1.5, false}},
			req:      req("CL", "2026-10-08", "2026-10-09", 2)},
		{name: "worked afternoon of this request breaks adjacency", existing: []limitRow{req("CL", "2026-10-12", "2026-10-13", 2)},
			req: limitRow{"E001", pend, "CL", "2026-10-08", "full", "2026-10-09", "first_half", 1.5, false}},
		{name: "half days that join count", existing: []limitRow{{"E001", pend, "CL", "2026-10-08", "second_half", "2026-10-09", "full", 1.5, false}},
			req: req("CL", "2026-10-12", "2026-10-12", 1), message: consec},
		{name: "chain before adds up", existing: []limitRow{req("PL", "2026-10-07", "2026-10-08", 2), req("PL", "2026-10-09", "2026-10-09", 1)},
			req: req("PL", "2026-10-12", "2026-10-14", 3), message: pl},
		{name: "chain after adds up", existing: []limitRow{req("PL", "2026-10-09", "2026-10-09", 1), req("PL", "2026-10-12", "2026-10-14", 3)},
			req: req("PL", "2026-10-07", "2026-10-08", 2), message: pl},
		{name: "PL 5 days with adjacent leave allowed", existing: []limitRow{req("PL", "2026-10-09", "2026-10-09", 1)},
			req: req("PL", "2026-10-12", "2026-10-15", 4)},
		{name: "max_consecutive_days 0 is no limit", req: func() limitRow {
			r := req("SL", "2026-10-12", "2026-10-19", 7)
			r.attach = true
			return r
		}()},

		// 7. quota types: days ≤ available
		{name: "CL 2 with 2 available allowed", existing: []limitRow{req("CL", "2026-10-20", "2026-10-21", 2)},
			req: req("CL", "2026-10-12", "2026-10-13", 2)},
		{name: "CL 2 with 1.5 available refused",
			existing: []limitRow{req("CL", "2026-10-20", "2026-10-21", 2), {"E001", pend, "CL", "2026-10-27", "first_half", "2026-10-27", "first_half", 0.5, false}},
			req:      req("CL", "2026-10-12", "2026-10-13", 2), message: "Not enough Casual Leave: 1.5 available, 2 needed."},
		{name: "CL used in the ledger refused", ledger: []adjustment{{"CL", -3}},
			req: req("CL", "2026-10-12", "2026-10-13", 2), message: "Not enough Casual Leave: 1 available, 2 needed."},
		{name: "EL 11 days with 10.5 available refused", req: req("EL", "2026-10-12", "2026-10-23", 11),
			message: "Not enough Earned Leave: 10.5 available, 11 needed."},
		{name: "no-quota BL needs no balance", req: req("BL", "2026-10-12", "2026-10-14", 3)},

		// 8. yearly cap: approved + pending of that type in the leave year + this request
		{name: "LOP up to the cap allowed",
			existing: []limitRow{mine(appr, "LOP", "2026-11-16", "2026-11-20", 5), mine(appr, "LOP", "2026-12-07", "2026-12-11", 5), req("LOP", "2027-01-04", "2027-01-07", 4)},
			req:      req("LOP", "2026-10-12", "2026-10-12", 1), warning: lopWarn},
		{name: "LOP beyond the cap refused",
			existing: []limitRow{mine(appr, "LOP", "2026-11-16", "2026-11-20", 5), mine(appr, "LOP", "2026-12-07", "2026-12-11", 5), req("LOP", "2027-01-04", "2027-01-07", 4)},
			req:      req("LOP", "2026-10-12", "2026-10-13", 2), message: lopCap},
		{name: "cancel_requested LOP counts toward the cap",
			existing: []limitRow{mine(appr, "LOP", "2026-11-16", "2026-11-20", 5), mine(appr, "LOP", "2026-12-07", "2026-12-11", 5), mine(approvals.CancelRequested, "LOP", "2027-01-04", "2027-01-07", 4)},
			req:      req("LOP", "2026-10-12", "2026-10-13", 2), message: lopCap},
		{name: "closed, other-year and other employees' LOP do not count",
			existing: []limitRow{mine(approvals.Rejected, "LOP", "2026-11-16", "2026-11-20", 5), mine(approvals.Cancelled, "LOP", "2026-12-07", "2026-12-11", 5),
				mine(appr, "LOP", "2026-03-16", "2026-03-20", 5), {"E002", appr, "LOP", "2027-01-04", "full", "2027-01-08", "full", 5, false},
				req("LOP", "2027-02-01", "2027-02-05", 5), req("LOP", "2027-02-08", "2027-02-12", 5)},
			req: req("LOP", "2026-10-12", "2026-10-16", 5), warning: lopWarn},

		// 9. times per employment
		{name: "second marriage leave refused", existing: []limitRow{mine(appr, "MRL", "2025-06-09", "2025-06-11", 3)},
			req: req("MRL", "2026-10-12", "2026-10-14", 3), message: mrlOnce},
		{name: "pending marriage leave counts", existing: []limitRow{req("MRL", "2026-12-07", "2026-12-09", 3)},
			req: req("MRL", "2026-10-12", "2026-10-14", 3), message: mrlOnce},
		{name: "rejected, cancelled and others' marriage leave do not count",
			existing: []limitRow{mine(approvals.Rejected, "MRL", "2025-06-09", "2025-06-11", 3), mine(approvals.Cancelled, "MRL", "2026-12-07", "2026-12-09", 3),
				{"E002", appr, "MRL", "2026-12-07", "full", "2026-12-09", "full", 3, false}},
			req: req("MRL", "2026-10-12", "2026-10-14", 3)},
		{name: "second of two allowed", rule: map[string]any{"max_times_per_employment": 2},
			existing: []limitRow{mine(appr, "BL", "2026-06-08", "2026-06-10", 3)}, req: req("BL", "2026-10-12", "2026-10-14", 3)},
		{name: "third of two refused", rule: map[string]any{"max_times_per_employment": 2},
			existing: []limitRow{mine(appr, "BL", "2026-06-08", "2026-06-10", 3), mine(appr, "BL", "2026-07-06", "2026-07-08", 3)},
			req:      req("BL", "2026-10-12", "2026-10-14", 3), message: "Bereavement Leave can be taken at most 2 times during employment."},

		// 10. attachment
		{name: "SL 3 days without attachment refused", req: req("SL", "2026-10-12", "2026-10-14", 3), message: sickNote},
		{name: "SL 2.5 days without attachment refused",
			req: limitRow{"E001", pend, "SL", "2026-10-12", "full", "2026-10-14", "first_half", 2.5, false}, message: sickNote},
		{name: "SL 3 days with attachment allowed", req: sl3},
		{name: "SL 2 days without attachment allowed", req: req("SL", "2026-10-12", "2026-10-13", 2)},

		// 11. LOP while paid leave is available: warning
		{name: "LOP with CL and EL available warns", req: req("LOP", "2026-10-12", "2026-10-12", 1), warning: lopWarn},
		{name: "LOP with only EL available warns", ledger: []adjustment{{"CL", -4}},
			req: req("LOP", "2026-10-12", "2026-10-12", 1), warning: "Loss of Pay applied for while paid leave is available (Earned Leave: 10.5)."},
		{name: "LOP with only CL available warns", ledger: []adjustment{{"EL", -10.5}},
			req: req("LOP", "2026-10-12", "2026-10-12", 1), warning: "Loss of Pay applied for while paid leave is available (Casual Leave: 4)."},
		{name: "pending CL reduces what the warning shows", existing: []limitRow{req("CL", "2026-10-20", "2026-10-20", 1)},
			req: req("LOP", "2026-10-12", "2026-10-12", 1), warning: "Loss of Pay applied for while paid leave is available (Casual Leave: 3, Earned Leave: 10.5)."},
		{name: "LOP with no paid leave has no warning", ledger: []adjustment{{"CL", -4}, {"EL", -10.5}},
			req: req("LOP", "2026-10-12", "2026-10-12", 1)},
		{name: "SL leaves only paid leave unwarned", req: req("SL", "2026-10-12", "2026-10-12", 1)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			app := testapp.New(t)
			for _, a := range c.ledger {
				addEntry(t, app, "E001", a.code, "2026-27", "adjustment", a.days)
			}
			if c.rule != nil {
				rule, err := leave.RuleFor(app, c.req.code, c.req.from)
				if err != nil {
					t.Fatal(err)
				}
				for k, v := range c.rule {
					rule.Set(k, v)
				}
				if err := app.Save(rule); err != nil {
					t.Fatal(err)
				}
			}
			for _, e := range c.existing {
				limitLeave(t, app, e)
			}
			lr := limitLeave(t, app, c.req)

			warnings, err := leave.LimitChecks(app, lr)
			if c.message != "" {
				var ue *approvals.UserError
				if !errors.As(err, &ue) || ue.Message != c.message {
					t.Fatalf("error = %v, want UserError %q", err, c.message)
				}
				return
			}
			if err != nil {
				t.Fatalf("refused: %v", err)
			}
			var want []string
			if c.warning != "" {
				want = []string{c.warning}
			}
			if len(warnings) != len(want) || (len(want) == 1 && warnings[0] != want[0]) {
				t.Errorf("warnings = %q, want %q", warnings, want)
			}
		})
	}
}
