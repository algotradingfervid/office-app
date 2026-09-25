package calendar

import (
	"slices"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// IsWorkingDay reports whether date (YYYY-MM-DD) is neither a weekly off nor a holiday.
func IsWorkingDay(app core.App, date string) (bool, error) {
	days, err := WorkingDays(app, date, date)
	return len(days) == 1, err
}

// WorkingDays returns the working dates from from to to (inclusive, YYYY-MM-DD), in order.
// It reads through app, so pass txApp inside a transaction.
func WorkingDays(app core.App, from, to string) ([]string, error) {
	start, err := time.Parse(time.DateOnly, from)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(time.DateOnly, to)
	if err != nil {
		return nil, err
	}

	setting, err := app.FindFirstRecordByData("settings", "key", "weekly_off")
	if err != nil {
		return nil, err
	}
	// weekday name -> weeks of the month that are off (1 = days 1-7, 2 = days 8-14, ...)
	var weeklyOff map[string][]int
	if err := setting.UnmarshalJSONField("value", &weeklyOff); err != nil {
		return nil, err
	}

	records, err := app.FindRecordsByFilter("holidays", "date >= {:from} && date <= {:to}", "", 0, 0,
		dbx.Params{"from": from, "to": to})
	if err != nil {
		return nil, err
	}
	holidays := map[string]bool{}
	for _, r := range records {
		holidays[r.GetString("date")] = true
	}

	var days []string
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		date := d.Format(time.DateOnly)
		week := (d.Day()-1)/7 + 1
		if holidays[date] || slices.Contains(weeklyOff[d.Weekday().String()], week) {
			continue
		}
		days = append(days, date)
	}
	return days, nil
}
