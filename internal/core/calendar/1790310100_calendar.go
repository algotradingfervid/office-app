package calendar

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Creates holidays and settings, and writes the weekly-off pattern so every database has it.
func init() {
	m.Register(func(app core.App) error {
		// API rules stay nil: no REST access.
		holidays := core.NewBaseCollection("holidays")
		holidays.Fields.Add(
			&core.TextField{Name: "date", Required: true, Pattern: `^\d{4}-\d{2}-\d{2}$`},
			&core.TextField{Name: "name", Required: true, Max: 100},
		)
		holidays.AddIndex("idx_holidays_date", true, "date", "")
		if err := app.Save(holidays); err != nil {
			return err
		}

		settings := core.NewBaseCollection("settings")
		settings.Fields.Add(
			&core.TextField{Name: "key", Required: true, Max: 100},
			&core.JSONField{Name: "value"},
		)
		settings.AddIndex("idx_settings_key", true, "key", "")
		if err := app.Save(settings); err != nil {
			return err
		}

		weeklyOff := core.NewRecord(settings)
		weeklyOff.Set("key", "weekly_off")
		weeklyOff.Set("value", map[string][]int{"Sunday": {1, 2, 3, 4, 5}, "Saturday": {2, 4}})
		return app.Save(weeklyOff)
	}, nil)
}
