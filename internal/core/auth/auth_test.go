package auth_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	"officeapp/internal/core/auth"
	"officeapp/internal/testapp"
)

// Data layer: the migration turns users into the employee directory, closed to the REST API.
func TestUsersCollectionShape(t *testing.T) {
	app := testapp.New(t)
	users, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"employee_code", "date_of_joining", "can_approve", "is_hr_admin", "active", "must_change_password"} {
		if users.Fields.GetByName(f) == nil {
			t.Errorf("users has no field %q", f)
		}
	}
	if users.ListRule != nil || users.ViewRule != nil || users.CreateRule != nil || users.UpdateRule != nil || users.DeleteRule != nil {
		t.Error("users API rules must all be nil (superuser only)")
	}
	dup := core.NewRecord(users)
	dup.Set("employee_code", "E001")
	dup.Set("name", "Duplicate")
	dup.SetPassword(auth.DemoPassword)
	if err := app.Save(dup); err == nil {
		t.Error("a second user with employee code E001 was saved")
	}
}

// Service layer: sign-in by employee code or e-mail; one message for every failure.
func TestLogin(t *testing.T) {
	app := testapp.New(t)
	inactive, _ := app.FindFirstRecordByData("users", "employee_code", "E003")
	inactive.Set("active", false)
	if err := app.Save(inactive); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name, identity, password string
		ok                       bool
	}{
		{"employee code", "E001", auth.DemoPassword, true},
		{"e-mail", "ravi@example.com", auth.DemoPassword, true},
		{"code with spaces", "  E002 ", auth.DemoPassword, true},
		{"wrong password", "E001", "wrong-password", false},
		{"unknown user", "E999", auth.DemoPassword, false},
		{"inactive user", "E003", auth.DemoPassword, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			token, err := auth.Login(app, c.identity, c.password)
			if c.ok && (err != nil || token == "") {
				t.Fatalf("want a token, got err %v", err)
			}
			if !c.ok && !errors.Is(err, auth.ErrBadLogin) {
				t.Fatalf("want ErrBadLogin, got %v", err)
			}
		})
	}
}

// Web layer: signed-out visitors go to /login; the cookie signs them in; cross-site posts are refused.
func TestSignInPages(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:           "home redirects signed-out visitors",
			Method:         http.MethodGet,
			URL:            "/",
			ExpectedStatus: http.StatusSeeOther,
			TestAppFactory: testapp.Factory,
		},
		{
			Name:            "login page renders with strict CSP",
			Method:          http.MethodGet,
			URL:             "/login",
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{"Employee code or e-mail", `action="/login"`},
			TestAppFactory:  testapp.Factory,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				if csp := res.Header.Get("Content-Security-Policy"); csp == "" {
					t.Error("no Content-Security-Policy header")
				}
			},
		},
		{
			Name:            "wrong password shows one plain message",
			Method:          http.MethodPost,
			URL:             "/login",
			Body:            testapp.Form(map[string]string{"identity": "E001", "password": "nope-nope-nope"}),
			Headers:         testapp.FormHeaders(""),
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedContent: []string{"Employee code, e-mail or password is incorrect."},
			TestAppFactory:  testapp.Factory,
		},
		{
			Name:           "right password sets an HttpOnly session cookie",
			Method:         http.MethodPost,
			URL:            "/login",
			Body:           testapp.Form(map[string]string{"identity": "E001", "password": auth.DemoPassword}),
			Headers:        testapp.FormHeaders(""),
			ExpectedStatus: http.StatusSeeOther,
			TestAppFactory: testapp.Factory,
			AfterTestFunc: func(t testing.TB, _ *tests.TestApp, res *http.Response) {
				for _, c := range res.Cookies() {
					if c.Name == auth.CookieName && c.Value != "" && c.HttpOnly {
						return
					}
				}
				t.Error("no HttpOnly session cookie set")
			},
		},
		{
			Name:            "cross-site login post is refused",
			Method:          http.MethodPost,
			URL:             "/login",
			Body:            testapp.Form(map[string]string{"identity": "E001", "password": auth.DemoPassword}),
			Headers:         map[string]string{"Content-Type": "application/x-www-form-urlencoded", "Sec-Fetch-Site": "cross-site"},
			ExpectedStatus:  http.StatusForbidden,
			ExpectedContent: []string{"Cross-origin request blocked."},
			TestAppFactory:  testapp.Factory,
		},
	}
	for _, s := range scenarios {
		s.Test(t)
	}
}

func TestHomeGreetsSignedInUser(t *testing.T) {
	app := testapp.New(t)
	cookie := testapp.SessionCookie(t, app, "E001")
	s := tests.ApiScenario{
		Method:          http.MethodGet,
		URL:             "/",
		Headers:         map[string]string{"Cookie": cookie},
		ExpectedStatus:  http.StatusOK,
		ExpectedContent: []string{"Welcome, Asha Rao", "Sign out"},
		TestAppFactory:  func(testing.TB) *tests.TestApp { return app },
	}
	s.Test(t)
}
