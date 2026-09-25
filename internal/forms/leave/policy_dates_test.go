package leave_test

import (
	"errors"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	"officeapp/internal/core/approvals"
	"officeapp/internal/core/clock"
	"officeapp/internal/forms/leave"
	"officeapp/internal/testapp"
)

// dateCheckToday is Monday 2026-10-05 in IST.
var dateCheckToday = clock.Fixed(time.Date(2026, 10, 5, 10, 0, 0, 0, clock.IST))

type leaveRow struct {
	emp, status, code         string
	from, fromSession, to, ts string
}

// saveLeave saves a request with its leave request, as SubmitRequest does before it calls Validate.
func saveLeave(t *testing.T, app *tests.TestApp, l leaveRow) *core.Record {
	t.Helper()
	emp, err := app.FindFirstRecordByData("users", "employee_code", l.emp)
	if err != nil {
		t.Fatal(err)
	}
	rule, err := leave.RuleFor(app, l.code, l.from)
	if err != nil {
		t.Fatal(err)
	}
	requests, err := app.FindCollectionByNameOrId("requests")
	if err != nil {
		t.Fatal(err)
	}
	req := core.NewRecord(requests)
	req.Set("form_type", "leave")
	req.Set("requester", emp.Id)
	req.Set("status", l.status)
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
	lr.Set("from_session", l.fromSession)
	lr.Set("to_date", l.to)
	lr.Set("to_session", l.ts)
	lr.Set("days", 1)
	lr.Set("leave_year", "2026-27")
	lr.Set("rule", rule.Id)
	if err := app.Save(lr); err != nil {
		t.Fatal(err)
	}
	return lr
}

func setProbationEnd(t *testing.T, app *tests.TestApp, code, date string) {
	t.Helper()
	u, err := app.FindFirstRecordByData("users", "employee_code", code)
	if err != nil {
		t.Fatal(err)
	}
	u.Set("probation_end", date)
	if err := app.Save(u); err != nil {
		t.Fatal(err)
	}
}

func pending(code, from, fs, to, ts string) leaveRow {
	return leaveRow{"E001", "pending", code, from, fs, to, ts}
}

func TestDateChecks(t *testing.T) {
	cases := []struct {
		name     string
		existing []leaveRow
		probEnd  string // E001's probation_end when set; the seed's 2024-12-03 is past
		req      leaveRow
		message  string // "" = allowed
		warning  string
	}{
		{name: "plain CL is allowed", req: pending("CL", "2026-10-12", "full", "2026-10-12", "full")},

		{name: "inactive type refused", req: pending("CO", "2026-10-12", "full", "2026-10-12", "full"),
			message: "Compensatory Off cannot be applied for at present."},
		{name: "EL in probation refused", probEnd: "2026-12-01", req: pending("EL", "2026-10-20", "full", "2026-10-20", "full"),
			message: "Earned Leave cannot be taken during probation."},
		{name: "EL on the last day of probation refused", probEnd: "2026-10-06", req: pending("EL", "2026-10-20", "full", "2026-10-20", "full"),
			message: "Earned Leave cannot be taken during probation."},
		{name: "EL on the day probation ends allowed", probEnd: "2026-10-05", req: pending("EL", "2026-10-20", "full", "2026-10-20", "full")},
		{name: "CL in probation allowed", probEnd: "2026-12-01", req: pending("CL", "2026-10-12", "full", "2026-10-12", "full")},

		{name: "end before start refused", req: pending("CL", "2026-10-13", "full", "2026-10-12", "full"),
			message: "The leave cannot end before it starts."},
		{name: "impossible date refused", req: pending("CL", "2026-11-31", "full", "2026-12-01", "full"),
			message: "Enter real dates."},
		{name: "crossing 31 March refused", req: pending("LOP", "2027-03-31", "full", "2027-04-01", "full"),
			message: "Leave cannot run past 31 March. Apply for the days up to 31 March and the days from 1 April as two requests."},
		{name: "ending on 31 March allowed", req: pending("LOP", "2027-03-29", "full", "2027-03-31", "full")},
		{name: "crossing 31 December allowed", req: pending("LOP", "2026-12-31", "full", "2027-01-01", "full")},

		{name: "overlap with own pending refused", existing: []leaveRow{pending("CL", "2026-10-12", "full", "2026-10-12", "full")},
			req: pending("SL", "2026-10-12", "full", "2026-10-12", "full"), message: "You already have leave applied for on some of these dates."},
		{name: "overlap with own approved refused", existing: []leaveRow{{"E001", "approved", "CL", "2026-10-12", "full", "2026-10-14", "full"}},
			req: pending("LOP", "2026-10-14", "first_half", "2026-10-14", "first_half"), message: "You already have leave applied for on some of these dates."},
		{name: "overlap with own cancel_requested refused", existing: []leaveRow{{"E001", "cancel_requested", "CL", "2026-10-12", "full", "2026-10-12", "full"}},
			req: pending("SL", "2026-10-09", "full", "2026-10-12", "first_half"), message: "You already have leave applied for on some of these dates."},
		{name: "cancelled request does not overlap", existing: []leaveRow{{"E001", "cancelled", "CL", "2026-10-12", "full", "2026-10-12", "full"}},
			req: pending("SL", "2026-10-12", "full", "2026-10-12", "full")},
		{name: "rejected request does not overlap", existing: []leaveRow{{"E001", "rejected", "CL", "2026-10-12", "full", "2026-10-12", "full"}},
			req: pending("SL", "2026-10-12", "full", "2026-10-12", "full")},
		{name: "another employee's leave does not overlap", existing: []leaveRow{{"E002", "approved", "CL", "2026-10-12", "full", "2026-10-12", "full"}},
			req: pending("SL", "2026-10-12", "full", "2026-10-12", "full")},
		{name: "first half then second half does not overlap", existing: []leaveRow{pending("CL", "2026-10-12", "first_half", "2026-10-12", "first_half")},
			req: pending("SL", "2026-10-12", "second_half", "2026-10-12", "second_half")},
		{name: "second half then first half does not overlap", existing: []leaveRow{pending("CL", "2026-10-12", "second_half", "2026-10-12", "second_half")},
			req: pending("SL", "2026-10-12", "first_half", "2026-10-12", "first_half")},
		{name: "same half overlaps", existing: []leaveRow{pending("CL", "2026-10-12", "second_half", "2026-10-12", "second_half")},
			req: pending("SL", "2026-10-12", "second_half", "2026-10-13", "full"), message: "You already have leave applied for on some of these dates."},
		{name: "ending first half before a second-half start does not overlap", existing: []leaveRow{pending("CL", "2026-10-13", "second_half", "2026-10-14", "full")},
			req: pending("SL", "2026-10-12", "full", "2026-10-13", "first_half")},
		{name: "full day overlaps a later request's second-half start", existing: []leaveRow{pending("CL", "2026-10-13", "second_half", "2026-10-14", "full")},
			req: pending("SL", "2026-10-12", "full", "2026-10-13", "full"), message: "You already have leave applied for on some of these dates."},
		{name: "day after an existing range does not overlap", existing: []leaveRow{pending("CL", "2026-10-12", "full", "2026-10-14", "full")},
			req: pending("SL", "2026-10-15", "full", "2026-10-15", "full")},
		{name: "day before an existing range does not overlap", existing: []leaveRow{pending("CL", "2026-10-12", "full", "2026-10-14", "full")},
			req: pending("SL", "2026-10-09", "full", "2026-10-09", "full")},
		{name: "range around an existing one overlaps", existing: []leaveRow{pending("CL", "2026-10-13", "full", "2026-10-13", "full")},
			req: pending("LOP", "2026-10-12", "full", "2026-10-14", "full"), message: "You already have leave applied for on some of these dates."},

		{name: "CL today allowed with short notice", req: pending("CL", "2026-10-05", "full", "2026-10-05", "full"),
			warning: "Short notice: Casual Leave asks for 1 day's notice."},
		{name: "CL yesterday refused (no backdating)", req: pending("CL", "2026-10-02", "full", "2026-10-02", "full"),
			message: "Casual Leave can start no earlier than 2026-10-05."},
		{name: "SL 7 days back allowed", req: pending("SL", "2026-09-28", "full", "2026-09-28", "full")},
		{name: "SL 8 days back refused", req: pending("SL", "2026-09-27", "full", "2026-09-28", "full"),
			message: "Sick Leave can start no earlier than 2026-09-28."},

		{name: "CL tomorrow has enough notice", req: pending("CL", "2026-10-06", "full", "2026-10-06", "full")},
		{name: "EL 7 days ahead has enough notice", req: pending("EL", "2026-10-12", "full", "2026-10-12", "full")},
		{name: "EL 6 days ahead is short notice", req: pending("EL", "2026-10-11", "full", "2026-10-12", "full"),
			warning: "Short notice: Earned Leave asks for 7 days' notice."},
		{name: "SL today needs no notice", req: pending("SL", "2026-10-05", "full", "2026-10-05", "full")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			app := testapp.New(t)
			if c.probEnd != "" {
				setProbationEnd(t, app, "E001", c.probEnd)
			}
			for _, e := range c.existing {
				saveLeave(t, app, e)
			}
			lr := saveLeave(t, app, c.req)

			warnings, err := leave.DateChecks(app, dateCheckToday, lr)
			if c.message == "" {
				if err != nil {
					t.Fatalf("refused: %v", err)
				}
			} else {
				var ue *approvals.UserError
				if !errors.As(err, &ue) || ue.Message != c.message {
					t.Fatalf("error = %v, want UserError %q", err, c.message)
				}
				return
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
