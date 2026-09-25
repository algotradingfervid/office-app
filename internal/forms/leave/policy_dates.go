package leave

import (
	"fmt"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/approvals"
	"officeapp/internal/core/clock"
)

// DateChecks runs checks 1–5 of design §5.3 on a saved leave_requests row: the type is active and
// allowed in probation, the dates are real and stay in one leave year, they do not overlap the
// employee's own pending, approved or cancel_requested leave, and the start is not backdated too far.
// A start inside the notice period is only a warning. Refusals are *approvals.UserError.
// The rule applied is the row's rule. It reads through app (txApp inside a transaction).
func DateChecks(app core.App, c clock.Clock, lr *core.Record) ([]string, error) {
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
	employee, err := app.FindRecordById("users", req.GetString("requester"))
	if err != nil {
		return nil, err
	}
	name := lt.GetString("name")
	today := clock.Today(c)

	if !lt.GetBool("active") {
		return nil, &approvals.UserError{Message: name + " cannot be applied for at present."}
	}
	if today < employee.GetString("probation_end") && !rule.GetBool("allowed_in_probation") {
		return nil, &approvals.UserError{Message: name + " cannot be taken during probation."}
	}

	from, to := lr.GetString("from_date"), lr.GetString("to_date")
	if _, err := time.Parse(time.DateOnly, from); err != nil {
		return nil, &approvals.UserError{Message: "Enter real dates."}
	}
	if _, err := time.Parse(time.DateOnly, to); err != nil {
		return nil, &approvals.UserError{Message: "Enter real dates."}
	}
	if to < from {
		return nil, &approvals.UserError{Message: "The leave cannot end before it starts."}
	}
	if LeaveYear(from) != LeaveYear(to) {
		return nil, &approvals.UserError{Message: "Leave cannot run past 31 March. " +
			"Apply for the days up to 31 March and the days from 1 April as two requests."}
	}

	overlapping, err := app.FindRecordsByFilter("leave_requests",
		"request.requester = {:emp} && request != {:req} && from_date <= {:to} && to_date >= {:from} && "+
			"(request.status = {:pending} || request.status = {:approved} || request.status = {:cancelReq})", "", 0, 0,
		dbx.Params{"emp": employee.Id, "req": req.Id, "from": from, "to": to, "pending": string(approvals.Pending),
			"approved": string(approvals.Approved), "cancelReq": string(approvals.CancelRequested)})
	if err != nil {
		return nil, err
	}
	start, end := halfDays(lr)
	for _, o := range overlapping {
		oStart, oEnd := halfDays(o)
		if start <= oEnd && oStart <= end {
			return nil, &approvals.UserError{Message: "You already have leave applied for on some of these dates."}
		}
	}

	now, _ := time.Parse(time.DateOnly, today)
	earliest := now.AddDate(0, 0, -rule.GetInt("max_backdate_days")).Format(time.DateOnly)
	if from < earliest {
		return nil, &approvals.UserError{Message: fmt.Sprintf("%s can start no earlier than %s.", name, earliest)}
	}

	var warnings []string
	// min_notice_days 0 means no notice is needed.
	if notice := rule.GetInt("min_notice_days"); notice > 0 && from < now.AddDate(0, 0, notice).Format(time.DateOnly) {
		days := fmt.Sprintf("%d days'", notice)
		if notice == 1 {
			days = "1 day's"
		}
		warnings = append(warnings, fmt.Sprintf("Short notice: %s asks for %s notice.", name, days))
	}
	return warnings, nil
}

// halfDays returns the first and last half day a leave request covers, as "YYYY-MM-DD1" (first half)
// or "YYYY-MM-DD2" (second half), so two requests overlap when each starts before the other ends.
func halfDays(lr *core.Record) (first, last string) {
	first, last = lr.GetString("from_date")+"1", lr.GetString("to_date")+"2"
	if lr.GetString("from_session") == "second_half" {
		first = lr.GetString("from_date") + "2"
	}
	if lr.GetString("to_session") == "first_half" {
		last = lr.GetString("to_date") + "1"
	}
	return first, last
}
