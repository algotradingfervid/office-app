package approvals_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	"officeapp/internal/core/approvals"
	"officeapp/internal/core/clock"
	"officeapp/internal/testapp"
)

// myRequestsApp has three requests, each of its own form type so each has its own summary:
// E001's older pending one (current approver E002), E001's newer approved one (no current approver),
// and E002's own request, which E001 must never see.
func myRequestsApp(t *testing.T) (app *tests.TestApp, older, newer *core.Record) {
	app = testapp.New(t)
	approvals.RegisterForm(app, "old", stubForm{"Casual leave, 12 Oct"})
	approvals.RegisterForm(app, "new", stubForm{"Sick leave, 20 Oct"})
	approvals.RegisterForm(app, "other", stubForm{"Not my request"})

	submit := func(formType, requester, approver string, at time.Time) *core.Record {
		req, err := approvals.SubmitRequest(app, clock.Fixed(at), approvals.SubmitInput{
			FormType: formType, FormVersion: 1,
			Requester: userID(t, app, requester), Approver: userID(t, app, approver),
			SaveForm: func(core.App, *core.Record) error { return nil },
		})
		if err != nil {
			t.Fatal(err)
		}
		return req
	}
	// 20:00 UTC on 1 October is already 2 October in India.
	older = submit("old", "E001", "E002", time.Date(2026, 10, 1, 20, 0, 0, 0, time.UTC))
	newer = submit("new", "E001", "E003", time.Date(2026, 10, 5, 9, 0, 0, 0, clock.IST))
	submit("other", "E002", "E003", time.Date(2026, 10, 6, 9, 0, 0, 0, clock.IST))

	newer.Set("status", string(approvals.Approved))
	newer.Set("current_approver", "")
	if err := app.Save(newer); err != nil {
		t.Fatal(err)
	}
	return app, older, newer
}

func TestMyRequestsListsOnlyMineNewestFirst(t *testing.T) {
	app, older, newer := myRequestsApp(t)
	s := tests.ApiScenario{
		Method:         http.MethodGet,
		URL:            "/requests",
		Headers:        map[string]string{"Cookie": testapp.SessionCookie(t, app, "E001")},
		ExpectedStatus: http.StatusOK,
		ExpectedContent: []string{
			"<h1>My requests</h1>",
			`<a href="/requests/` + older.Id + `">Casual leave, 12 Oct</a>`,
			"2 Oct 2026", "Pending", "Ravi Kumar",
			`<a href="/requests/` + newer.Id + `">Sick leave, 20 Oct</a>`,
			"5 Oct 2026", "Approved",
		},
		NotExpectedContent: []string{"Not my request", "1 Oct 2026", "Meena Iyer", "You have not made any requests yet."},
		TestAppFactory:     func(testing.TB) *tests.TestApp { return app },
		AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
			body := readBody(t, res)
			if strings.Index(body, "Sick leave") > strings.Index(body, "Casual leave") {
				t.Error("the newer request is not listed first")
			}
		},
	}
	s.Test(t)
}

func TestMyRequestsShowsTheOtherEmployeesOwnRequestsToThem(t *testing.T) {
	app, _, _ := myRequestsApp(t)
	s := tests.ApiScenario{
		Method:             http.MethodGet,
		URL:                "/requests",
		Headers:            map[string]string{"Cookie": testapp.SessionCookie(t, app, "E002")},
		ExpectedStatus:     http.StatusOK,
		ExpectedContent:    []string{"Not my request", "6 Oct 2026", "Meena Iyer"},
		NotExpectedContent: []string{"Casual leave", "Sick leave"},
		TestAppFactory:     func(testing.TB) *tests.TestApp { return app },
	}
	s.Test(t)
}

func TestMyRequestsStatusLabels(t *testing.T) {
	cases := map[approvals.Status]string{
		approvals.Pending:         "<td>Pending</td>",
		approvals.Approved:        "<td>Approved</td>",
		approvals.Rejected:        "<td>Rejected</td>",
		approvals.Cancelled:       "<td>Cancelled</td>",
		approvals.CancelRequested: "<td>Cancellation requested</td>",
	}
	for status, want := range cases {
		t.Run(string(status), func(t *testing.T) {
			app, older, _ := myRequestsApp(t)
			older.Set("status", string(status))
			if err := app.Save(older); err != nil {
				t.Fatal(err)
			}
			s := tests.ApiScenario{
				Method:          http.MethodGet,
				URL:             "/requests",
				Headers:         map[string]string{"Cookie": testapp.SessionCookie(t, app, "E001")},
				ExpectedStatus:  http.StatusOK,
				ExpectedContent: []string{want},
				TestAppFactory:  func(testing.TB) *tests.TestApp { return app },
			}
			s.Test(t)
		})
	}
}

func TestMyRequestsPageSignedOutAndEmpty(t *testing.T) {
	app := testapp.New(t)
	scenarios := []tests.ApiScenario{
		{
			Name:           "signed-out visitors go to sign in",
			Method:         http.MethodGet,
			URL:            "/requests",
			ExpectedStatus: http.StatusSeeOther,
			TestAppFactory: testapp.Factory,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				if loc := res.Header.Get("Location"); loc != "/login" {
					t.Errorf("Location = %q, want /login", loc)
				}
			},
		},
		{
			Name:               "an employee with no requests sees the empty state",
			Method:             http.MethodGet,
			URL:                "/requests",
			Headers:            map[string]string{"Cookie": testapp.SessionCookie(t, app, "E001")},
			ExpectedStatus:     http.StatusOK,
			ExpectedContent:    []string{"<h1>My requests</h1>", "You have not made any requests yet."},
			NotExpectedContent: []string{"<table"},
			TestAppFactory:     func(testing.TB) *tests.TestApp { return app },
		},
	}
	for _, s := range scenarios {
		s.Test(t)
	}
}

func readBody(t testing.TB, res *http.Response) string {
	t.Helper()
	b := new(strings.Builder)
	if _, err := io.Copy(b, res.Body); err != nil {
		t.Fatal(err)
	}
	return b.String()
}
