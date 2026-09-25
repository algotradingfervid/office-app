package leave

import (
	"fmt"
	"slices"
	"time"

	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/approvals"
	"officeapp/internal/core/calendar"
)

// Days counts a leave request's days (design §5.3); a second_half start or a first_half end counts 0.5.
// countMode "working_days" skips weekly offs and holidays and the request must start and end on a
// working day. "calendar_days" (maternity) counts every date, Sundays and holidays included; only a
// half day on a non-working day is refused. Refusals are *approvals.UserError.
// It reads through app (txApp inside a transaction).
func Days(app core.App, countMode, from, fromSession, to, toSession string) (float64, error) {
	start, err := time.Parse(time.DateOnly, from)
	if err != nil {
		return 0, err
	}
	end, err := time.Parse(time.DateOnly, to)
	if err != nil {
		return 0, err
	}
	sessions := []string{"full", "first_half", "second_half"}
	switch {
	case !slices.Contains(sessions, fromSession) || !slices.Contains(sessions, toSession):
		return 0, &approvals.UserError{Message: "Choose a session: full day, first half or second half."}
	case end.Before(start):
		return 0, &approvals.UserError{Message: "The leave cannot end before it starts."}
	case from == to && fromSession != toSession:
		return 0, &approvals.UserError{Message: "A one-day leave has a single session: full day, first half or second half."}
	case from != to && fromSession == "first_half":
		return 0, &approvals.UserError{Message: "A leave of more than one day can start with the second half of a day, not the first half."}
	case from != to && toSession == "second_half":
		return 0, &approvals.UserError{Message: "A leave of more than one day can end with the first half of a day, not the second half."}
	}

	var days float64
	switch countMode {
	case "working_days":
		working, err := calendar.WorkingDays(app, from, to)
		if err != nil {
			return 0, err
		}
		if len(working) == 0 || working[0] != from || working[len(working)-1] != to {
			return 0, &approvals.UserError{Message: "Leave must start and end on a working day."}
		}
		days = float64(len(working))
	case "calendar_days":
		for _, d := range [][2]string{{from, fromSession}, {to, toSession}} {
			if d[1] == "full" {
				continue
			}
			working, err := calendar.IsWorkingDay(app, d[0])
			if err != nil {
				return 0, err
			}
			if !working {
				return 0, &approvals.UserError{Message: "A half day must be on a working day."}
			}
		}
		days = end.Sub(start).Hours()/24 + 1
	default:
		return 0, fmt.Errorf("unknown count mode %q", countMode)
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
