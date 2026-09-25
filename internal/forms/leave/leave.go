// Package leave is the Leave form: leave types, versioned rules, the balance ledger,
// policy checks and the apply page. It uses core modules only through their root packages.
package leave

import "github.com/pocketbase/pocketbase/core"

// parts are added by the module's own files from their init(), so no story edits a shared
// registration list: a file that adds routes, hooks or seed data calls addPart in its init().
var parts []func(app core.App)

func addPart(p func(app core.App)) { parts = append(parts, p) }

// Register runs every part of the module (in file-name order, as Go runs init functions).
func Register(app core.App) {
	for _, p := range parts {
		p(app)
	}
}
