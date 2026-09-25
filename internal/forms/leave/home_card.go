package leave

import (
	"embed"
	"html/template"
	"strings"

	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/clock"
	"officeapp/internal/core/home"
)

//go:embed templates/balance_card.html
var balanceCardFS embed.FS

var balanceCardTmpl = template.Must(template.ParseFS(balanceCardFS, "templates/balance_card.html"))

func init() {
	addPart(func(app core.App) {
		home.AddCard(app, balanceCard(clock.System{}))
	})
}

// balanceCard lists the signed-in employee's balance and available days of each active quota type
// for the current leave year.
func balanceCard(c clock.Clock) home.Card {
	return func(e *core.RequestEvent) (template.HTML, error) {
		year := LeaveYear(clock.Today(c))
		types, err := AvailableBalances(e.App, e.Auth.Id, year)
		if err != nil {
			return "", err
		}
		var b strings.Builder
		if err := balanceCardTmpl.Execute(&b, map[string]any{"Year": year, "Types": types}); err != nil {
			return "", err
		}
		return template.HTML(b.String()), nil
	}
}
