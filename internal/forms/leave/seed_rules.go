package leave

import "github.com/pocketbase/pocketbase/core"

// defaultRule is one leave type with its rule from 2026-04-01 (design §5.2, leave policy §3–§5).
// What 0 means differs per field; see the leave_rules fields in 1790310300_leave.go.
type defaultRule struct {
	code, name, credit, count           string
	perYear, monthly, maxRun, yearlyCap float64
	maxTimes, notice, backdate          int
	halfDay, encashable, probation      bool
	attachAfter, carryCap               float64
	active                              bool
}

var defaultRules = []defaultRule{
	{"CL", "Casual Leave", "monthly", "working_days", 12, 1, 2, 0, 0, 1, 0, true, false, true, 0, 0, true},
	{"SL", "Sick Leave", "yearly", "working_days", 8, 0, 0, 0, 0, 0, 7, true, false, true, 2, 0, true},
	{"EL", "Earned Leave", "monthly", "working_days", 18, 1.5, 0, 0, 0, 7, 0, true, true, false, 0, 30, true},
	{"ML", "Maternity Leave", "per_event", "calendar_days", 0, 0, 182, 0, 0, 0, 0, false, false, true, 0, 0, true},
	{"PL", "Paternity Leave", "per_event", "working_days", 0, 0, 5, 0, 0, 0, 0, false, false, true, 0, 0, true},
	{"BL", "Bereavement Leave", "per_event", "working_days", 0, 0, 3, 0, 0, 0, 0, false, false, true, 0, 0, true},
	{"MRL", "Marriage Leave", "per_event", "working_days", 0, 0, 3, 0, 1, 0, 0, false, false, true, 0, 0, true},
	{"CO", "Compensatory Off", "per_event", "working_days", 0, 0, 0, 0, 0, 0, 0, false, false, true, 0, 0, false},
	{"LOP", "Loss of Pay", "per_event", "working_days", 0, 0, 5, 15, 0, 0, 0, false, false, true, 0, 0, true},
}

// seedRules writes the nine leave types and their first rules. It runs in the migration,
// so every database (not only demo ones) starts with the policy defaults.
func seedRules(app core.App, types, rules *core.Collection) error {
	for _, d := range defaultRules {
		lt := core.NewRecord(types)
		lt.Set("code", d.code)
		lt.Set("name", d.name)
		lt.Set("active", d.active)
		if err := app.Save(lt); err != nil {
			return err
		}
		r := core.NewRecord(rules)
		r.Set("leave_type", lt.Id)
		r.Set("effective_from", "2026-04-01")
		r.Set("days_per_year", d.perYear)
		r.Set("credit_method", d.credit)
		r.Set("monthly_credit", d.monthly)
		r.Set("count_mode", d.count)
		r.Set("max_consecutive_days", d.maxRun)
		r.Set("yearly_cap", d.yearlyCap)
		r.Set("max_times_per_employment", d.maxTimes)
		r.Set("min_notice_days", d.notice)
		r.Set("max_backdate_days", d.backdate)
		r.Set("half_day_allowed", d.halfDay)
		r.Set("attachment_after_days", d.attachAfter)
		r.Set("carry_forward_cap", d.carryCap)
		r.Set("encashable", d.encashable)
		r.Set("allowed_in_probation", d.probation)
		if err := app.Save(r); err != nil {
			return err
		}
	}
	return nil
}
