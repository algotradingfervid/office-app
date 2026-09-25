package leave

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/approvals"
)

// TypeAvailable is one leave type's balance and the part of it not held by pending requests.
type TypeAvailable struct {
	TypeBalance
	Available float64
}

// Available is the balance minus the days of that leave type and leave year in the employee's pending
// requests (design §5.4). Only pending requests hold days: an approved one is already debited in the ledger.
func Available(app core.App, employeeID, typeCode, leaveYear string) (float64, error) {
	b, err := Balance(app, employeeID, typeCode, leaveYear)
	if err != nil {
		return 0, err
	}
	p, err := pendingDays(app, employeeID, typeCode, leaveYear)
	if err != nil {
		return 0, err
	}
	return b - p, nil
}

// AvailableBalances is Balances with each type's available days.
func AvailableBalances(app core.App, employeeID, leaveYear string) ([]TypeAvailable, error) {
	balances, err := Balances(app, employeeID, leaveYear)
	if err != nil {
		return nil, err
	}
	out := make([]TypeAvailable, 0, len(balances))
	for _, b := range balances {
		p, err := pendingDays(app, employeeID, b.Code, leaveYear)
		if err != nil {
			return nil, err
		}
		out = append(out, TypeAvailable{TypeBalance: b, Available: b.Balance - p})
	}
	return out, nil
}

func pendingDays(app core.App, employeeID, typeCode, leaveYear string) (float64, error) {
	rows, err := app.FindRecordsByFilter("leave_requests",
		"request.requester = {:emp} && request.status = {:status} && leave_type.code = {:code} && leave_year = {:year}",
		"", 0, 0,
		dbx.Params{"emp": employeeID, "status": string(approvals.Pending), "code": typeCode, "year": leaveYear})
	if err != nil {
		return 0, err
	}
	sum := 0.0
	for _, r := range rows {
		sum += r.GetFloat("days")
	}
	return sum, nil
}
