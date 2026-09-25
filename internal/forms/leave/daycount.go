package leave

import (
	"time"

	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/approvals"
	"officeapp/internal/core/calendar"
)

// Days counts a leave request's days (design §5.3): countMode "working_days" skips weekly offs and
// holidays, "calendar_days" counts every date; a second_half start or a first_half end counts 0.5.
// The request must start and end on a working day in both modes, which also refuses a half day on a
// non-working day. Refusals are *approvals.UserError. It reads through app (txApp inside a transaction).
func Days(app core.App, countMode, from, fromSession, to, toSession string) (float64, error) {
	start, err := time.Parse(time.DateOnly, from)
	if err != nil {
		return 0, err
	}
	end, err := time.Parse(time.DateOnly, to)
	if err != nil {
		return 0, err
	}
	switch {
	case end.Before(start):
		return 0, &approvals.UserError{Message: "The leave cannot end before it starts."}
	case from == to && fromSession != toSession:
		return 0, &approvals.UserError{Message: "A one-day leave has a single session: full day, first half or second half."}
	case from != to && fromSession == "first_half":
		return 0, &approvals.UserError{Message: "A leave of more than one day can start with the second half of a day, not the first half."}
	case from != to && toSession == "second_half":
		return 0, &approvals.UserError{Message: "A leave of more than one day can end with the first half of a day, not the second half."}
	}

	working, err := calendar.WorkingDays(app, from, to)
	if err != nil {
		return 0, err
	}
	if len(working) == 0 || working[0] != from || working[len(working)-1] != to {
		return 0, &approvals.UserError{Message: "Leave must start and end on a working day."}
	}

	days := float64(len(working))
	if countMode == "calendar_days" {
		days = end.Sub(start).Hours()/24 + 1
	}
	if fromSession == "second_half" {
		days -= 0.5
	}
	// A one-day half leave has the same session in both fields, so only one of these applies.
	if toSession == "first_half" {
		days -= 0.5
	}
	return days, nil
}
