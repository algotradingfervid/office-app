package calendar

import (
	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/seed"
)

func init() {
	addPart(func(app core.App) {
		seed.Add(app, "calendar.holidays", seedHolidays)
	})
}

// seedHolidays adds the national holidays of leave year 2026-27 and Diwali.
func seedHolidays(txApp core.App) error {
	holidays, err := txApp.FindCollectionByNameOrId("holidays")
	if err != nil {
		return err
	}
	demo := []struct{ date, name string }{
		{"2026-08-15", "Independence Day"},
		{"2026-10-02", "Gandhi Jayanti"},
		{"2026-11-09", "Diwali"},
		{"2027-01-26", "Republic Day"},
	}
	for _, d := range demo {
		r := core.NewRecord(holidays)
		r.Set("date", d.date)
		r.Set("name", d.name)
		if err := txApp.Save(r); err != nil {
			return err
		}
	}
	return nil
}
