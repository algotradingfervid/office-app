package leave_test

import (
	"errors"
	"testing"

	"officeapp/internal/core/approvals"
	"officeapp/internal/forms/leave"
	"officeapp/internal/testapp"
)

func TestDays(t *testing.T) {
	app := testapp.New(t) // weekly offs: Sundays, 2nd and 4th Saturdays; holiday 2026-10-02
	cases := []struct {
		name              string
		mode              string
		from, fromSession string
		to, toSession     string
		want              float64
	}{
		{"two full days", "working_days", "2026-10-12", "full", "2026-10-13", "full", 2},
		{"weekend skipped", "working_days", "2026-10-09", "full", "2026-10-12", "full", 2},
		{"holiday skipped, 1st Saturday counted", "working_days", "2026-10-01", "full", "2026-10-05", "full", 3},
		{"one full day", "working_days", "2026-10-12", "full", "2026-10-12", "full", 1},
		{"first half", "working_days", "2026-10-12", "first_half", "2026-10-12", "first_half", 0.5},
		{"second half", "working_days", "2026-10-12", "second_half", "2026-10-12", "second_half", 0.5},
		{"second half to first half", "working_days", "2026-10-12", "second_half", "2026-10-13", "first_half", 1},
		{"second half to full", "working_days", "2026-10-12", "second_half", "2026-10-14", "full", 2.5},
		{"full to first half", "working_days", "2026-10-12", "full", "2026-10-14", "first_half", 2.5},
		{"maternity calendar days", "calendar_days", "2026-10-01", "full", "2026-10-31", "full", 31},
		{"calendar days across a year", "calendar_days", "2026-12-31", "full", "2027-01-01", "full", 2},
		{"calendar days half", "calendar_days", "2026-10-09", "second_half", "2026-10-12", "full", 3.5},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := leave.Days(app, c.mode, c.from, c.fromSession, c.to, c.toSession)
			if err != nil {
				t.Fatal(err)
			}
			if got != c.want {
				t.Errorf("Days = %v, want %v", got, c.want)
			}
		})
	}
}

func TestDaysRefused(t *testing.T) {
	app := testapp.New(t)
	cases := []struct {
		name              string
		mode              string
		from, fromSession string
		to, toSession     string
		message           string
	}{
		{"starts on a Sunday", "working_days", "2026-10-11", "full", "2026-10-13", "full",
			"Leave must start and end on a working day."},
		{"ends on a 2nd Saturday", "working_days", "2026-10-09", "full", "2026-10-10", "full",
			"Leave must start and end on a working day."},
		{"half day on a holiday", "working_days", "2026-10-02", "first_half", "2026-10-02", "first_half",
			"Leave must start and end on a working day."},
		{"calendar days starting on a Sunday", "calendar_days", "2026-10-11", "full", "2026-10-31", "full",
			"Leave must start and end on a working day."},
		{"calendar days ending on a holiday", "calendar_days", "2026-09-01", "full", "2026-10-02", "full",
			"Leave must start and end on a working day."},
		{"to before from", "working_days", "2026-10-13", "full", "2026-10-12", "full",
			"The leave cannot end before it starts."},
		{"one day, two sessions", "working_days", "2026-10-12", "first_half", "2026-10-12", "second_half",
			"A one-day leave has a single session: full day, first half or second half."},
		{"starts with a first half", "working_days", "2026-10-12", "first_half", "2026-10-13", "full",
			"A leave of more than one day can start with the second half of a day, not the first half."},
		{"ends with a second half", "working_days", "2026-10-12", "full", "2026-10-13", "second_half",
			"A leave of more than one day can end with the first half of a day, not the second half."},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := leave.Days(app, c.mode, c.from, c.fromSession, c.to, c.toSession)
			var userErr *approvals.UserError
			if !errors.As(err, &userErr) {
				t.Fatalf("err = %v, want a UserError", err)
			}
			if userErr.Message != c.message {
				t.Errorf("message = %q, want %q", userErr.Message, c.message)
			}
		})
	}
}
