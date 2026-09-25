// Package calendar owns the holiday list and the weekly-off pattern, and answers
// "is this date a working day?" for every form.
package calendar

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
