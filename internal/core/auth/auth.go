// Package auth signs employees in and out with a cookie that holds a PocketBase auth token,
// and loads the signed-in user into e.Auth on every request.
package auth

import (
	"embed"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/seed"
	"officeapp/internal/core/web"
)

// CookieName holds the session token.
const CookieName = "oa_session"

// ErrBadLogin is returned for every failed sign-in, whatever the reason.
var ErrBadLogin = errors.New("auth: bad login")

// badLoginMessage is the only message a failed sign-in ever shows.
const badLoginMessage = "Employee code, e-mail or password is incorrect."

//go:embed templates
var templatesFS embed.FS

var pages = web.NewPages(templatesFS)

// Register adds the session middleware, the sign-in routes and the demo users.
func Register(app core.App) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.BindFunc(loadSession)
		se.Router.GET("/login", showLogin)
		se.Router.POST("/login", doLogin)
		se.Router.POST("/logout", doLogout)
		return se.Next()
	})
	seed.Add(app, "auth", seedUsers)
}

// RequireUser sends signed-out visitors to the sign-in page. Bind it on a module's routes.
func RequireUser(e *core.RequestEvent) error {
	if e.Auth == nil || e.Auth.Collection().Name != "users" {
		return e.Redirect(http.StatusSeeOther, "/login")
	}
	return e.Next()
}

// Login checks an identity (employee code or e-mail) and password and returns a session token.
func Login(app core.App, identity, password string) (string, error) {
	identity = strings.TrimSpace(identity)
	field := "employee_code"
	if strings.Contains(identity, "@") {
		field = "email"
	}
	user, err := app.FindFirstRecordByData("users", field, identity)
	if err != nil || !user.GetBool("active") || !user.ValidatePassword(password) {
		return "", ErrBadLogin
	}
	return user.NewAuthToken()
}

func loadSession(e *core.RequestEvent) error {
	if c, err := e.Request.Cookie(CookieName); err == nil && e.Auth == nil {
		user, err := e.App.FindAuthRecordByToken(c.Value, core.TokenTypeAuth)
		if err == nil && user.Collection().Name == "users" && user.GetBool("active") {
			e.Auth = user
		}
	}
	return e.Next()
}

func showLogin(e *core.RequestEvent) error {
	if e.Auth != nil {
		return e.Redirect(http.StatusSeeOther, "/")
	}
	return pages.Render(e, "login", "Sign in", loginForm{})
}

type loginForm struct {
	Identity string
	Error    string
}

func doLogin(e *core.RequestEvent) error {
	f := loginForm{Identity: e.Request.FormValue("identity")}
	token, err := Login(e.App, f.Identity, e.Request.FormValue("password"))
	if err != nil {
		if !errors.Is(err, ErrBadLogin) {
			return err
		}
		f.Error = badLoginMessage
		return pages.RenderStatus(e, http.StatusUnauthorized, "login", "Sign in", f)
	}
	users, err := e.App.FindCachedCollectionByNameOrId("users")
	if err != nil {
		return err
	}
	e.SetCookie(sessionCookie(e, token, int(users.AuthToken.Duration)))
	return e.Redirect(http.StatusSeeOther, "/")
}

func doLogout(e *core.RequestEvent) error {
	e.SetCookie(sessionCookie(e, "", -1))
	return e.Redirect(http.StatusSeeOther, "/login")
}

// sessionCookie is Secure everywhere except plain-HTTP previews on this machine.
func sessionCookie(e *core.RequestEvent, value string, maxAge int) *http.Cookie {
	host, _, err := net.SplitHostPort(e.Request.Host)
	if err != nil {
		host = e.Request.Host
	}
	return &http.Cookie{
		Name:     CookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   host != "127.0.0.1" && host != "localhost",
		SameSite: http.SameSiteLaxMode,
	}
}
