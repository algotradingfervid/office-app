package leave

import (
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/approvals"
	"officeapp/internal/core/calendar"
)

// LimitChecks runs checks 6–11 of design §5.3 on a saved, pending leave_requests row: the consecutive-day
// limit (adjacent requests of the same type added together), the available balance of quota types, the
// yearly cap, the times-per-employment limit and the attachment. LOP while CL or EL is available is only a
// warning. The days are counted again with Days; other requests count when pending, approved or
// cancel_requested. A rule number of 0 means no limit. Refusals are *approvals.UserError.
// It reads through app (txApp inside a transaction).
func LimitChecks(app core.App, lr *core.Record) ([]string, error) {
	lt, err := app.FindRecordById("leave_types", lr.GetString("leave_type"))
	if err != nil {
		return nil, err
	}
	rule, err := app.FindRecordById("leave_rules", lr.GetString("rule"))
	if err != nil {
		return nil, err
	}
	req, err := app.FindRecordById("requests", lr.GetString("request"))
	if err != nil {
		return nil, err
	}
	name, code, employee, year := lt.GetString("name"), lt.GetString("code"), req.GetString("requester"), lr.GetString("leave_year")

	// Balance reads 0 for a type it cannot use, so refuse an inactive type before reading balances.
	if !lt.GetBool("active") {
		return nil, &approvals.UserError{Message: name + " cannot be applied for at present."}
	}
	days, err := Days(app, rule.GetString("count_mode"), lr.GetString("from_date"), lr.GetString("from_session"),
		lr.GetString("to_date"), lr.GetString("to_session"))
	if err != nil {
		return nil, err
	}

	// The employee's other requests of this type, in every leave year, that hold or have used days.
	others, err := app.FindRecordsByFilter("leave_requests",
		"request.requester = {:emp} && request != {:req} && leave_type = {:type} && "+
			"(request.status = {:pending} || request.status = {:approved} || request.status = {:cancelReq})", "", 0, 0,
		dbx.Params{"emp": employee, "req": req.Id, "type": lt.Id, "pending": string(approvals.Pending),
			"approved": string(approvals.Approved), "cancelReq": string(approvals.CancelRequested)})
	if err != nil {
		return nil, err
	}

	if limit := rule.GetFloat("max_consecutive_days"); limit > 0 {
		adjacent, err := adjacentDays(app, lr, others)
		if err != nil {
			return nil, err
		}
		if days+adjacent > limit {
			return nil, &approvals.UserError{Message: fmt.Sprintf("%s allows at most %g consecutive days.", name, limit)}
		}
	}

	if rule.GetFloat("days_per_year") > 0 {
		available, err := Available(app, employee, code, year)
		if err != nil {
			return nil, err
		}
		// Available already holds this pending request's stored days; give them back.
		available += lr.GetFloat("days")
		if days > available {
			return nil, &approvals.UserError{Message: fmt.Sprintf("Not enough %s: %g available, %g needed.", name, available, days)}
		}
	}

	if limit := rule.GetFloat("yearly_cap"); limit > 0 {
		used := days
		for _, o := range others {
			if o.GetString("leave_year") == year {
				used += o.GetFloat("days")
			}
		}
		if used > limit {
			return nil, &approvals.UserError{Message: fmt.Sprintf("%s is limited to %g days in a leave year. Contact HR.", name, limit)}
		}
	}

	if limit := rule.GetInt("max_times_per_employment"); limit > 0 && len(others) >= limit {
		times := "only once"
		if limit > 1 {
			times = fmt.Sprintf("at most %d times", limit)
		}
		return nil, &approvals.UserError{Message: fmt.Sprintf("%s can be taken %s during employment.", name, times)}
	}

	if after := rule.GetFloat("attachment_after_days"); after > 0 && days > after && lr.GetString("attachment") == "" {
		return nil, &approvals.UserError{Message: fmt.Sprintf("%s of more than %g days needs an attachment.", name, after)}
	}

	if code != "LOP" {
		return nil, nil
	}
	balances, err := AvailableBalances(app, employee, year)
	if err != nil {
		return nil, err
	}
	var paid []string
	for _, b := range balances {
		if (b.Code == "CL" || b.Code == "EL") && b.Available > 0 {
			paid = append(paid, fmt.Sprintf("%s: %g", b.Name, b.Available))
		}
	}
	if len(paid) == 0 {
		return nil, nil
	}
	return []string{"Loss of Pay applied for while paid leave is available (" + strings.Join(paid, ", ") + ")."}, nil
}

// adjacentDays adds up the days of the requests in others that join lr's run of leave, directly or through
// each other. The run only grows outwards, so no request is counted twice.
func adjacentDays(app core.App, lr *core.Record, others []*core.Record) (float64, error) {
	first, last, sum := lr, lr, 0.0
	for grown := true; grown; {
		grown = false
		for _, o := range others {
			before, err := joins(app, o, first)
			if err != nil {
				return 0, err
			}
			after, err := joins(app, last, o)
			if err != nil {
				return 0, err
			}
			switch {
			case before:
				first = o
			case after:
				last = o
			default:
				continue
			}
			sum += o.GetFloat("days")
			grown = true
		}
	}
	return sum, nil
}

// joins reports whether leave b starts right after leave a: only non-working days between them and no
// worked half day (a ending at the first half, or b starting at the second half).
func joins(app core.App, a, b *core.Record) (bool, error) {
	end, start := a.GetString("to_date"), b.GetString("from_date")
	if end >= start || a.GetString("to_session") == "first_half" || b.GetString("from_session") == "second_half" {
		return false, nil
	}
	working, err := calendar.WorkingDays(app, end, start)
	if err != nil {
		return false, err
	}
	for _, d := range working {
		if d > end && d < start {
			return false, nil
		}
	}
	return true, nil
}
