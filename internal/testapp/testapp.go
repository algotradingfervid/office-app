// Package testapp builds a throwaway app with every module registered and the demo data seeded.
package testapp

import (
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/pocketbase/pocketbase/tests"

	"officeapp/internal/core/auth"
	"officeapp/internal/core/seed"
	"officeapp/internal/modules"
)

var (
	templateOnce sync.Once
	templateDir  string
	templateErr  error
)

// New returns a migrated, seeded app on its own copy of the database. It is cleaned up with the test.
// The demo data is seeded once per test binary and cloned, because hashing passwords is slow.
func New(t testing.TB) *tests.TestApp {
	t.Helper()
	templateOnce.Do(buildTemplate)
	if templateErr != nil {
		t.Fatal(templateErr)
	}
	app, err := tests.NewTestApp(templateDir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(app.Cleanup)
	modules.Register(app, nil)
	return app
}

func buildTemplate() {
	templateDir, templateErr = os.MkdirTemp("", "officeapp-testdata-")
	if templateErr != nil {
		return
	}
	app, err := tests.NewTestApp(templateDir)
	if err != nil {
		templateErr = err
		return
	}
	modules.Register(app, nil)
	templateErr = seed.Run(app)
	if templateErr == nil {
		templateErr = app.ClearBootstrap() // flush and close the DB before copying it
	}
	if templateErr == nil {
		templateErr = os.RemoveAll(templateDir)
		if templateErr == nil {
			templateErr = os.Rename(app.DataDir(), templateDir)
		}
	}
}

// Factory adapts New for tests.ApiScenario.TestAppFactory.
func Factory(t testing.TB) *tests.TestApp { return New(t) }

// SessionCookie signs in a demo user by employee code and returns the Cookie header value.
func SessionCookie(t testing.TB, app *tests.TestApp, code string) string {
	t.Helper()
	token, err := auth.Login(app, code, auth.DemoPassword)
	if err != nil {
		t.Fatal(err)
	}
	return auth.CookieName + "=" + token
}

// Form encodes fields as an application/x-www-form-urlencoded body.
func Form(fields map[string]string) *strings.Reader {
	v := url.Values{}
	for k, val := range fields {
		v.Set(k, val)
	}
	return strings.NewReader(v.Encode())
}

// FormHeaders are the headers a same-origin browser form post sends.
func FormHeaders(cookie string) map[string]string {
	h := map[string]string{
		"Content-Type":   "application/x-www-form-urlencoded",
		"Sec-Fetch-Site": "same-origin",
	}
	if cookie != "" {
		h["Cookie"] = cookie
	}
	return h
}
