package approvals_test

import (
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

// inboxApp has three requests, each of its own form type so each has its own summary.
// The page uses the real clock, so submission times are relative to now:
// E001's request to E002 from 3 days ago, E003's request to E002 from 1 day ago,
// and E001's request to E003 from today, which E002 must never see.
func inboxApp(t *testing.T) (app *tests.TestApp, older, newer *core.Record) {
	app = testapp.New(t)
	approvals.RegisterForm(app, "old", stubForm{"Casual leave, 12 Oct"})
	approvals.RegisterForm(app, "new", stubForm{"Sick leave, 20 Oct"})
	approvals.RegisterForm(app, "other", stubForm{"Waiting for Meena"})

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
	now := time.Now()
	older = submit("old", "E001", "E002", now.AddDate(0, 0, -3))
	newer = submit("new", "E003", "E002", now.AddDate(0, 0, -1))
	submit("other", "E001", "E003", now)
	return app, older, newer
}

func istDate(t time.Time) string { return t.In(clock.IST).Format("2 Jan 2006") }

func TestInboxListsOnlyWhatWaitsForMeOldestFirst(t *testing.T) {
	app, older, newer := inboxApp(t)
	now := time.Now()
	s := tests.ApiScenario{
		Method:         http.MethodGet,
		URL:            "/approvals",
		Headers:        map[string]string{"Cookie": testapp.SessionCookie(t, app, "E002")},
		ExpectedStatus: http.StatusOK,
		ExpectedContent: []string{
			"<h1>Pending my approval</h1>",
			`<a href="/requests/` + older.Id + `">Casual leave, 12 Oct</a>`,
			"<td>Asha Rao</td>", "<td>" + istDate(now.AddDate(0, 0, -3)) + "</td>", "<td>3 days</td>",
			`<a href="/requests/` + newer.Id + `">Sick leave, 20 Oct</a>`,
			"<td>Meena Iyer</td>", "<td>" + istDate(now.AddDate(0, 0, -1)) + "</td>", "<td>1 day</td>",
		},
		NotExpectedContent: []string{"Waiting for Meena", "Nothing is waiting for you."},
		TestAppFactory:     func(testing.TB) *tests.TestApp { return app },
		AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
			body := readBody(t, res)
			if strings.Index(body, "Casual leave") > strings.Index(body, "Sick leave") {
				t.Error("the older request is not listed first")
			}
		},
	}
	s.Test(t)
}

func TestInboxShowsTheOtherApproverOnlyTheirs(t *testing.T) {
	app, _, _ := inboxApp(t)
	s := tests.ApiScenario{
		Method:             http.MethodGet,
		URL:                "/approvals",
		Headers:            map[string]string{"Cookie": testapp.SessionCookie(t, app, "E003")},
		ExpectedStatus:     http.StatusOK,
		ExpectedContent:    []string{"Waiting for Meena", "<td>Asha Rao</td>", "<td>today</td>"},
		NotExpectedContent: []string{"Casual leave", "Sick leave"},
		TestAppFactory:     func(testing.TB) *tests.TestApp { return app },
	}
	s.Test(t)
}

// Only pending and cancel_requested requests wait for an approver: a rejected request keeps its
// current_approver, and the other statuses must not show up either.
func TestInboxShowsOnlyRequestsThatWaitForAction(t *testing.T) {
	cases := map[approvals.Status]string{
		approvals.Pending:         "<td>Pending</td>",
		approvals.CancelRequested: "<td>Cancellation requested</td>",
		approvals.Approved:        "",
		approvals.Rejected:        "",
		approvals.Cancelled:       "",
	}
	for status, want := range cases {
		t.Run(string(status), func(t *testing.T) {
			app, older, _ := inboxApp(t)
			older.Set("status", string(status))
			if err := app.Save(older); err != nil {
				t.Fatal(err)
			}
			s := tests.ApiScenario{
				Method:          http.MethodGet,
				URL:             "/approvals",
				Headers:         map[string]string{"Cookie": testapp.SessionCookie(t, app, "E002")},
				ExpectedStatus:  http.StatusOK,
				ExpectedContent: []string{"Sick leave"},
				TestAppFactory:  func(testing.TB) *tests.TestApp { return app },
			}
			if want == "" {
				s.NotExpectedContent = []string{"Casual leave"}
			} else {
				s.ExpectedContent = append(s.ExpectedContent, "Casual leave", want)
			}
			s.Test(t)
		})
	}
}

func TestInboxRefusesNonApproversAndSignedOut(t *testing.T) {
	refused, empty := testapp.New(t), testapp.New(t) // an app serves one scenario only
	scenarios := []tests.ApiScenario{
		{
			Name:            "an employee who is not an approver is refused",
			Method:          http.MethodGet,
			URL:             "/approvals",
			Headers:         map[string]string{"Cookie": testapp.SessionCookie(t, refused, "E001")},
			ExpectedStatus:  http.StatusForbidden,
			ExpectedContent: []string{"Only approvers can see this page."},
			TestAppFactory:  func(testing.TB) *tests.TestApp { return refused },
		},
		{
			Name:           "signed-out visitors go to sign in",
			Method:         http.MethodGet,
			URL:            "/approvals",
			ExpectedStatus: http.StatusSeeOther,
			TestAppFactory: testapp.Factory,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				if loc := res.Header.Get("Location"); loc != "/login" {
					t.Errorf("Location = %q, want /login", loc)
				}
			},
		},
		{
			Name:               "an approver with nothing waiting sees the empty state",
			Method:             http.MethodGet,
			URL:                "/approvals",
			Headers:            map[string]string{"Cookie": testapp.SessionCookie(t, empty, "E002")},
			ExpectedStatus:     http.StatusOK,
			ExpectedContent:    []string{"<h1>Pending my approval</h1>", "Nothing is waiting for you."},
			NotExpectedContent: []string{"<table"},
			TestAppFactory:     func(testing.TB) *tests.TestApp { return empty },
		},
	}
	for _, s := range scenarios {
		s.Test(t)
	}
}
