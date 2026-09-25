// Package approvals owns requests of every form, their approval steps and the status rules
// (submit, forward, approve as final, reject, cancel). Forms plug in through its hooks interface.
package approvals

import "github.com/pocketbase/pocketbase/core"

// parts are added by the module's own files from their init(), so no story edits a shared
// registration list: a file that adds routes, hooks or seed data calls addPart in its init().
var parts []func(app core.App)

func addPart(p func(app core.App)) { parts = append(parts, p) } //nolint:unused // called by feature files as they are built

// Register runs every part of the module (in file-name order, as Go runs init functions).
func Register(app core.App) {
	for _, p := range parts {
		p(app)
	}
}
