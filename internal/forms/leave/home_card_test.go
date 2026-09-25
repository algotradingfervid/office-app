package leave_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	"officeapp/internal/core/approvals"
	"officeapp/internal/core/clock"
	"officeapp/internal/forms/leave"
	"officeapp/internal/testapp"
)

// renderBalanceCard renders the card for employee code at the given instant.
func renderBalanceCard(t *testing.T, app *tests.TestApp, code string, at time.Time) string {
	t.Helper()
	user, err := app.FindFirstRecordByData("users", "employee_code", code)
	if err != nil {
		t.Fatal(err)
	}
	html, err := leave.BalanceCardAt(clock.Fixed(at))(&core.RequestEvent{App: app, Auth: user})
	if err != nil {
		t.Fatal(err)
	}
	return string(html)
}

func TestBalanceCardShowsBalanceAndAvailableForTheLeaveYear(t *testing.T) {
	app := testapp.New(t)
	requestLeave(t, app, "E001", approvals.Pending, "CL", "2026-27", 1)
	got := renderBalanceCard(t, app, "E001", time.Date(2026, 10, 1, 9, 0, 0, 0, clock.IST))
	for _, want := range []string{
		"Leave balances 2026-27",
		"<td>Casual Leave</td><td>4</td><td>3</td>",
		"<td>Earned Leave</td><td>10.5</td><td>10.5</td>",
		"<td>Sick Leave</td><td>8</td><td>8</td>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("card does not contain %q:\n%s", want, got)
		}
	}
}

func TestBalanceCardUsesTheClocksLeaveYear(t *testing.T) {
	app := testapp.New(t)
	// 20:00 UTC on 31 March 2027 is already 1 April in India: leave year 2027-28 has no entries yet.
	got := renderBalanceCard(t, app, "E001", time.Date(2027, 3, 31, 20, 0, 0, 0, time.UTC))
	for _, want := range []string{"Leave balances 2027-28", "<td>Casual Leave</td><td>0</td><td>0</td>"} {
		if !strings.Contains(got, want) {
			t.Errorf("card does not contain %q:\n%s", want, got)
		}
	}
}

func TestBalanceCardIsOnTheHomePage(t *testing.T) {
	app := testapp.New(t)
	s := tests.ApiScenario{
		Method:          http.MethodGet,
		URL:             "/",
		Headers:         map[string]string{"Cookie": testapp.SessionCookie(t, app, "E001")},
		ExpectedStatus:  http.StatusOK,
		ExpectedContent: []string{"Leave balances", "<td>Casual Leave</td>", "<td>Earned Leave</td>", "<td>Sick Leave</td>"},
		TestAppFactory:  func(testing.TB) *tests.TestApp { return app },
	}
	s.Test(t)
}
