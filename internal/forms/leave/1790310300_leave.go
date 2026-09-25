package leave

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Creates the leave collections (design §5.2) and seeds the default types and rules.
// API rules stay nil: no REST access. PocketBase numbers cannot be empty, so each leave_rules
// number says below what 0 means.
// leave_requests.request and leave_ledger.request (relations to requests) come in their own migration.
func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		datePattern := `^\d{4}-\d{2}-\d{2}$`
		yearPattern := `^\d{4}-\d{2}$`
		zero := 0.0

		types := core.NewBaseCollection("leave_types")
		types.Fields.Add(
			&core.TextField{Name: "code", Required: true, Max: 10},
			&core.TextField{Name: "name", Required: true, Max: 100},
			&core.BoolField{Name: "active"},
		)
		types.AddIndex("idx_leave_types_code", true, "code", "")
		if err := app.Save(types); err != nil {
			return err
		}

		rules := core.NewBaseCollection("leave_rules")
		rules.Fields.Add(
			&core.RelationField{Name: "leave_type", CollectionId: types.Id, MaxSelect: 1, Required: true},
			&core.TextField{Name: "effective_from", Required: true, Pattern: datePattern},
			&core.NumberField{Name: "days_per_year", Min: &zero}, // 0 = no quota (no balance)
			&core.SelectField{Name: "credit_method", Values: []string{"yearly", "monthly", "per_event"}, MaxSelect: 1, Required: true},
			&core.NumberField{Name: "monthly_credit", Min: &zero},
			&core.SelectField{Name: "count_mode", Values: []string{"working_days", "calendar_days"}, MaxSelect: 1, Required: true},
			&core.NumberField{Name: "max_consecutive_days", Min: &zero},                    // 0 = no limit
			&core.NumberField{Name: "yearly_cap", Min: &zero},                              // 0 = no cap
			&core.NumberField{Name: "max_times_per_employment", Min: &zero, OnlyInt: true}, // 0 = no limit
			&core.NumberField{Name: "min_notice_days", Min: &zero, OnlyInt: true},          // 0 = no notice needed
			&core.NumberField{Name: "max_backdate_days", Min: &zero, OnlyInt: true},        // 0 = a real zero: no backdating
			&core.BoolField{Name: "half_day_allowed"},
			&core.NumberField{Name: "attachment_after_days", Min: &zero}, // 0 = no attachment needed
			&core.NumberField{Name: "carry_forward_cap", Min: &zero},     // 0 = a real zero: everything lapses at year close
			&core.BoolField{Name: "encashable"},
			&core.BoolField{Name: "allowed_in_probation"},
		)
		rules.AddIndex("idx_leave_rules_type_from", true, "leave_type, effective_from", "")
		if err := app.Save(rules); err != nil {
			return err
		}

		sessions := []string{"full", "first_half", "second_half"}
		requests := core.NewBaseCollection("leave_requests")
		requests.Fields.Add(
			&core.RelationField{Name: "leave_type", CollectionId: types.Id, MaxSelect: 1, Required: true},
			&core.TextField{Name: "from_date", Required: true, Pattern: datePattern},
			&core.SelectField{Name: "from_session", Values: sessions, MaxSelect: 1, Required: true},
			&core.TextField{Name: "to_date", Required: true, Pattern: datePattern},
			&core.SelectField{Name: "to_session", Values: sessions, MaxSelect: 1, Required: true},
			&core.NumberField{Name: "days", Required: true, Min: &zero},
			&core.TextField{Name: "leave_year", Required: true, Pattern: yearPattern},
			&core.TextField{Name: "reason", Max: 500},
			&core.FileField{Name: "attachment", MaxSelect: 1, MaxSize: 5 << 20, Protected: true,
				MimeTypes: []string{"application/pdf", "image/jpeg", "image/png"}},
			&core.RelationField{Name: "rule", CollectionId: rules.Id, MaxSelect: 1, Required: true},
		)
		if err := app.Save(requests); err != nil {
			return err
		}

		ledger := core.NewBaseCollection("leave_ledger")
		ledger.Fields.Add(
			&core.RelationField{Name: "employee", CollectionId: users.Id, MaxSelect: 1, Required: true},
			&core.RelationField{Name: "leave_type", CollectionId: types.Id, MaxSelect: 1, Required: true},
			&core.TextField{Name: "leave_year", Required: true, Pattern: yearPattern},
			&core.SelectField{Name: "entry_type", MaxSelect: 1, Required: true,
				Values: []string{"opening", "credit", "debit", "reversal", "carry_forward", "lapse", "adjustment"}},
			&core.NumberField{Name: "days", Required: true},
			&core.TextField{Name: "period_key", Max: 20},
			&core.TextField{Name: "note", Max: 500},
			&core.RelationField{Name: "created_by", CollectionId: users.Id, MaxSelect: 1},
		)
		ledger.AddIndex("idx_leave_ledger_period", true, "employee, leave_type, entry_type, leave_year, period_key", "period_key != ''")
		if err := app.Save(ledger); err != nil {
			return err
		}

		return seedRules(app, types, rules)
	}, nil)
}
