package auth

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Shapes PocketBase's built-in users collection into the employee directory.
func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		// The browser never talks to the PocketBase API; the app's own handlers authorise every access.
		users.ListRule, users.ViewRule, users.CreateRule, users.UpdateRule, users.DeleteRule = nil, nil, nil, nil, nil

		users.Fields.GetByName("email").(*core.EmailField).Required = false
		users.Fields.GetByName("password").(*core.PasswordField).Min = 10
		users.Fields.GetByName("name").(*core.TextField).Required = true
		datePattern := `^\d{4}-\d{2}-\d{2}$`
		users.Fields.Add(
			&core.TextField{Name: "employee_code", Required: true, Max: 20},
			&core.TextField{Name: "department", Max: 100},
			&core.TextField{Name: "designation", Max: 100},
			&core.TextField{Name: "date_of_joining", Pattern: datePattern},
			&core.TextField{Name: "probation_end", Pattern: datePattern},
			&core.BoolField{Name: "can_approve"},
			&core.BoolField{Name: "is_hr_admin"},
			&core.BoolField{Name: "active"},
			&core.BoolField{Name: "must_change_password"},
		)
		users.AddIndex("idx_users_employee_code", true, "employee_code", "")
		users.PasswordAuth.IdentityFields = []string{"email", "employee_code"}
		users.AuthToken.Duration = 7 * 24 * 60 * 60
		return app.Save(users)
	}, nil)
}
