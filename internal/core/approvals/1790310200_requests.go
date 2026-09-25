package approvals

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Creates requests (one row per request of any form) and approval_steps (its append-only history).
// API rules stay nil: no REST access, the app's own handlers authorise everything.
func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		requests := core.NewBaseCollection("requests")
		requests.Fields.Add(
			&core.TextField{Name: "form_type", Required: true, Max: 50},
			&core.RelationField{Name: "requester", CollectionId: users.Id, Required: true, MaxSelect: 1},
			&core.SelectField{Name: "status", Required: true, MaxSelect: 1,
				Values: []string{"pending", "approved", "rejected", "cancelled", "cancel_requested"}},
			&core.RelationField{Name: "current_approver", CollectionId: users.Id, MaxSelect: 1},
			&core.RelationField{Name: "final_approver", CollectionId: users.Id, MaxSelect: 1},
			&core.DateField{Name: "submitted_at", Required: true},
			&core.NumberField{Name: "form_version", Required: true, OnlyInt: true},
			&core.BoolField{Name: "recorded_by_hr"},
		)
		if err := app.Save(requests); err != nil {
			return err
		}

		steps := core.NewBaseCollection("approval_steps")
		steps.Fields.Add(
			&core.RelationField{Name: "request", CollectionId: requests.Id, Required: true, MaxSelect: 1},
			&core.NumberField{Name: "seq", Required: true, OnlyInt: true},
			&core.RelationField{Name: "actor", CollectionId: users.Id, Required: true, MaxSelect: 1},
			&core.SelectField{Name: "action", Required: true, MaxSelect: 1,
				Values: []string{"submitted", "forwarded", "approved_final", "rejected", "cancelled", "cancel_requested",
					"cancel_approved", "cancel_declined", "reassigned", "recorded_by_hr"}},
			&core.TextField{Name: "comment", Max: 1000},
			&core.RelationField{Name: "to_user", CollectionId: users.Id, MaxSelect: 1},
			&core.JSONField{Name: "warnings"},
			&core.AutodateField{Name: "created", OnCreate: true},
		)
		steps.AddIndex("idx_approval_steps_request_seq", true, "request, seq", "")
		return app.Save(steps)
	}, nil)
}
