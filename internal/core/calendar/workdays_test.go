package calendar_test

import (
	"slices"
	"testing"

	"officeapp/internal/core/calendar"
	"officeapp/internal/testapp"
)

func TestIsWorkingDay(t *testing.T) {
	app := testapp.New(t)
	cases := []struct {
		date    string
		working bool
	}{
		{"2026-10-10", false}, // 2nd Saturday
		{"2026-10-17", true},  // 3rd Saturday
		{"2026-10-24", false}, // 4th Saturday
		{"2026-10-31", true},  // 5th Saturday
		{"2026-10-11", false}, // Sunday
		{"2026-10-02", false}, // Gandhi Jayanti (seeded holiday)
		{"2026-10-12", true},  // Monday
		{"2026-11-07", true},  // 1st Saturday (day 7, last day of week 1)
		{"2026-11-14", false}, // 2nd Saturday (day 14, last day of week 2)
		{"2026-11-28", false}, // 4th Saturday (day 28, last day of week 4)
		{"2026-11-09", false}, // Diwali (seeded holiday)
		{"2027-01-26", false}, // Republic Day (seeded holiday)
		{"2026-08-15", false}, // Independence Day (seeded holiday, also a 3rd Saturday)
	}
	for _, c := range cases {
		t.Run(c.date, func(t *testing.T) {
			got, err := calendar.IsWorkingDay(app, c.date)
			if err != nil {
				t.Fatal(err)
			}
			if got != c.working {
				t.Errorf("IsWorkingDay(%s) = %v, want %v", c.date, got, c.working)
			}
		})
	}
}

func TestWorkingDays(t *testing.T) {
	app := testapp.New(t)
	got, err := calendar.WorkingDays(app, "2026-10-09", "2026-10-13")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"2026-10-09", "2026-10-12", "2026-10-13"}
	if !slices.Equal(got, want) {
		t.Errorf("WorkingDays = %v, want %v", got, want)
	}

	got, err = calendar.WorkingDays(app, "2026-10-13", "2026-10-09")
	if err != nil || len(got) != 0 {
		t.Errorf("WorkingDays(from > to) = %v, %v; want empty, nil", got, err)
	}
}

func TestMalformedDate(t *testing.T) {
	app := testapp.New(t)
	for _, d := range []string{"", "2026-10-32", "10/10/2026", "2026-1-5"} {
		if _, err := calendar.IsWorkingDay(app, d); err == nil {
			t.Errorf("IsWorkingDay(%q) returned no error", d)
		}
		if _, err := calendar.WorkingDays(app, "2026-10-01", d); err == nil {
			t.Errorf("WorkingDays(to %q) returned no error", d)
		}
		if _, err := calendar.WorkingDays(app, d, "2026-10-01"); err == nil {
			t.Errorf("WorkingDays(from %q) returned no error", d)
		}
	}
}

func TestMissingWeeklyOff(t *testing.T) {
	app := testapp.New(t)
	rec, err := app.FindFirstRecordByData("settings", "key", "weekly_off")
	if err != nil {
		t.Fatal(err)
	}
	if err := app.Delete(rec); err != nil {
		t.Fatal(err)
	}
	if _, err := calendar.IsWorkingDay(app, "2026-10-12"); err == nil {
		t.Error("IsWorkingDay without weekly_off returned no error")
	}
}

func TestCollectionsClosedToAPI(t *testing.T) {
	app := testapp.New(t)
	for _, name := range []string{"holidays", "settings"} {
		c, err := app.FindCollectionByNameOrId(name)
		if err != nil {
			t.Fatal(err)
		}
		if c.ListRule != nil || c.ViewRule != nil || c.CreateRule != nil || c.UpdateRule != nil || c.DeleteRule != nil {
			t.Errorf("%s API rules must all be nil", name)
		}
	}
}
