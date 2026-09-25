package leave

import (
	"fmt"
	"strconv"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// LeaveYear returns the leave year (1 April – 31 March) of a YYYY-MM-DD date, written "2026-27".
// date must already be a valid YYYY-MM-DD (every date field is pattern-validated).
func LeaveYear(date string) string {
	year, _ := strconv.Atoi(date[:4])
	if date[5:7] < "04" {
		year--
	}
	return fmt.Sprintf("%d-%02d", year, (year+1)%100)
}

// RuleFor returns the rule of a leave type (by code) in effect on date:
// the one with the latest effective_from on or before date.
func RuleFor(app core.App, typeCode, date string) (*core.Record, error) {
	found, err := app.FindRecordsByFilter("leave_rules",
		"leave_type.code = {:code} && effective_from <= {:date}", "-effective_from", 1, 0,
		dbx.Params{"code": typeCode, "date": date})
	if err != nil {
		return nil, err
	}
	if len(found) == 0 {
		return nil, fmt.Errorf("no %s leave rule in effect on %s", typeCode, date)
	}
	return found[0], nil
}
