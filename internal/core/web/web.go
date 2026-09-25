// Package web is the shared page kit: the layout, static assets, rendering and the
// security middleware every page gets. Modules render their own templates through Pages.
package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// AppName is the working name shown in the layout until /name-it picks one.
const AppName = "Office App"

//go:embed static
var staticFS embed.FS

//go:embed templates
var layoutFS embed.FS

var crossOrigin = http.NewCrossOriginProtection()

// Register serves /static/ and adds the security middleware to every route.
func Register(app core.App) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		static, _ := fs.Sub(staticFS, "static")
		se.Router.GET("/static/{path...}", apis.Static(static, false))
		se.Router.BindFunc(securityHeaders)
		se.Router.BindFunc(rejectCrossOrigin)
		return se.Next()
	})
}

// securityHeaders sets a strict CSP: scripts and styles only from our own origin, no inline code.
// The PocketBase dashboard (/_/) is left alone: it is reachable only over Tailscale and needs its own inline code.
func securityHeaders(e *core.RequestEvent) error {
	if strings.HasPrefix(e.Request.URL.Path, "/_/") {
		return e.Next()
	}
	h := e.Response.Header()
	h.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; frame-ancestors 'none'; form-action 'self'")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Referrer-Policy", "same-origin")
	return e.Next()
}

// rejectCrossOrigin blocks state-changing requests that come from another site (CSRF).
func rejectCrossOrigin(e *core.RequestEvent) error {
	if err := crossOrigin.Check(e.Request); err != nil {
		return e.ForbiddenError("Cross-origin request blocked.", nil)
	}
	return e.Next()
}
