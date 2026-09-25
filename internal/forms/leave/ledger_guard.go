package leave

import (
	"errors"
	"fmt"
	"math"

	"github.com/pocketbase/pocketbase/core"
)

// halfDayFields are the number fields that count days, per collection: each must be a multiple of 0.5.
var halfDayFields = map[string][]string{
	"leave_ledger":   {"days"},
	"leave_requests": {"days"},
	"leave_rules":    {"days_per_year", "monthly_credit", "max_consecutive_days", "yearly_cap", "attachment_after_days", "carry_forward_cap"},
}

var errLedgerAppendOnly = errors.New("leave ledger entries are never changed or deleted; add a correcting entry")

// The ledger is append-only (design §5): these hooks also stop edits from the superuser dashboard.
func init() {
	addPart(func(app core.App) {
		app.OnRecordUpdate("leave_ledger").BindFunc(func(e *core.RecordEvent) error { return errLedgerAppendOnly })
		app.OnRecordDelete("leave_ledger").BindFunc(func(e *core.RecordEvent) error { return errLedgerAppendOnly })
		app.OnRecordValidate("leave_ledger", "leave_requests", "leave_rules").BindFunc(func(e *core.RecordEvent) error {
			for _, f := range halfDayFields[e.Record.Collection().Name] {
				if v := e.Record.GetFloat(f); v*2 != math.Trunc(v*2) {
					return fmt.Errorf("%s must be a multiple of 0.5, got %v", f, v)
				}
			}
			return e.Next()
		})
	})
}
