// Package home is the signed-in landing page. Form modules add cards to it with AddCard,
// so home never imports a form.
package home

import (
	"embed"
	"html/template"

	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/auth"
	"officeapp/internal/core/web"
)

// Card renders one block of the home page for the signed-in user.
type Card func(e *core.RequestEvent) (template.HTML, error)

const storeKey = "home.cards"

//go:embed templates
var templatesFS embed.FS

var pages = web.NewPages(templatesFS)

// AddCard appends a card to the home page. Call it from a module's Register.
func AddCard(app core.App, c Card) {
	list, _ := app.Store().Get(storeKey).([]Card)
	app.Store().Set(storeKey, append(list, c))
}

// Register adds GET / for signed-in users.
func Register(app core.App) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/{$}", show).BindFunc(auth.RequireUser)
		return se.Next()
	})
}

func show(e *core.RequestEvent) error {
	cards, _ := e.App.Store().Get(storeKey).([]Card)
	html := make([]template.HTML, 0, len(cards))
	for _, c := range cards {
		h, err := c(e)
		if err != nil {
			return err
		}
		html = append(html, h)
	}
	return pages.Render(e, "home", "Home", html)
}
