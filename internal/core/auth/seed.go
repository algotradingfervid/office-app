package auth

import "github.com/pocketbase/pocketbase/core"

// DemoPassword is the password of every demo user (development and previews only).
const DemoPassword = "demo-pass-2026"

// seedUsers creates one employee, two approvers and one HR admin.
func seedUsers(txApp core.App) error {
	users, err := txApp.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}
	demo := []struct {
		code, name, email     string
		canApprove, isHRAdmin bool
	}{
		{"E001", "Asha Rao", "asha@example.com", false, false},
		{"E002", "Ravi Kumar", "ravi@example.com", true, false},
		{"E003", "Meena Iyer", "meena@example.com", true, false},
		{"E004", "Farah Khan", "", false, true},
	}
	for _, d := range demo {
		r := core.NewRecord(users)
		r.Set("employee_code", d.code)
		r.Set("name", d.name)
		r.Set("email", d.email)
		r.Set("department", "Operations")
		r.Set("date_of_joining", "2024-06-03")
		r.Set("probation_end", "2024-12-03")
		r.Set("can_approve", d.canApprove)
		r.Set("is_hr_admin", d.isHRAdmin)
		r.Set("active", true)
		r.SetPassword(DemoPassword)
		if err := txApp.Save(r); err != nil {
			return err
		}
	}
	return nil
}
