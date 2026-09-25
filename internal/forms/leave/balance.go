package leave

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// TypeBalance is one leave type's balance for an employee and leave year.
type TypeBalance struct {
	Code, Name string
	Balance    float64
}

// Balance is the sum of the employee's ledger entries for a leave type (by code) and leave year (design §5.4).
// employeeID is the users record id. Days are multiples of 0.5, so the float sum is exact.
func Balance(app core.App, employeeID, typeCode, leaveYear string) (float64, error) {
	entries, err := app.FindRecordsByFilter("leave_ledger",
		"employee = {:emp} && leave_type.code = {:code} && leave_year = {:year}", "", 0, 0,
		dbx.Params{"emp": employeeID, "code": typeCode, "year": leaveYear})
	if err != nil {
		return 0, err
	}
	sum := 0.0
	for _, e := range entries {
		sum += e.GetFloat("days")
	}
	return sum, nil
}

// Balances returns the balance of every active quota type, ordered by code.
// A quota type is one with a rule version whose days_per_year is not 0 (0 = no quota).
func Balances(app core.App, employeeID, leaveYear string) ([]TypeBalance, error) {
	types, err := app.FindRecordsByFilter("leave_types",
		"active = true && leave_rules_via_leave_type.days_per_year ?> 0", "code", 0, 0)
	if err != nil {
		return nil, err
	}
	out := make([]TypeBalance, 0, len(types))
	for _, lt := range types {
		b, err := Balance(app, employeeID, lt.GetString("code"), leaveYear)
		if err != nil {
			return nil, err
		}
		out = append(out, TypeBalance{Code: lt.GetString("code"), Name: lt.GetString("name"), Balance: b})
	}
	return out, nil
}
