package leave

import (
	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/seed"
)

func init() {
	addPart(func(app core.App) {
		seed.Add(app, "leave.balances", seedBalances)
	})
}

// seedBalances gives each demo employee opening balances for 2026-27 (design §5.4: HR loads them at go-live).
func seedBalances(txApp core.App) error {
	ledger, err := txApp.FindCollectionByNameOrId("leave_ledger")
	if err != nil {
		return err
	}
	opening := []struct {
		code string
		days float64
	}{{"CL", 4}, {"SL", 8}, {"EL", 10.5}}
	for _, emp := range []string{"E001", "E002", "E003", "E004"} {
		user, err := txApp.FindFirstRecordByData("users", "employee_code", emp)
		if err != nil {
			return err
		}
		for _, o := range opening {
			lt, err := txApp.FindFirstRecordByData("leave_types", "code", o.code)
			if err != nil {
				return err
			}
			r := core.NewRecord(ledger)
			r.Set("employee", user.Id)
			r.Set("leave_type", lt.Id)
			r.Set("leave_year", "2026-27")
			r.Set("entry_type", "opening")
			r.Set("days", o.days)
			r.Set("note", "Demo opening balance")
			if err := txApp.Save(r); err != nil {
				return err
			}
		}
	}
	return nil
}
