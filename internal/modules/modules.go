// Package modules is the one list of the app's modules, in dependency order.
// Adding a module (a new form) adds one line here; nothing else is shared.
package modules

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"

	"officeapp/internal/core/approvals"
	"officeapp/internal/core/auth"
	"officeapp/internal/core/calendar"
	"officeapp/internal/core/home"
	"officeapp/internal/core/seed"
	"officeapp/internal/core/web"
	"officeapp/internal/forms/leave"
)

// Register wires every module into the app. Tests call it on a test app.
func Register(app core.App, root *cobra.Command) {
	web.Register(app)
	auth.Register(app)
	calendar.Register(app)
	approvals.Register(app)
	leave.Register(app)
	home.Register(app)

	if root != nil {
		root.AddCommand(seed.Command(app))
	}
}
