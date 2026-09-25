package approvals

import (
	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/clock"
)

// Form is what a form module plugs into approvals (the form hooks, design spec §4).
// Every hook that takes txApp runs inside the action's transaction and must read and write through it.
type Form interface {
	// Validate runs on submit (and preview) and again on approve as final. An error blocks the action;
	// a *UserError is shown to the user as-is, any other error is internal. Warnings are shown to the approver.
	Validate(txApp core.App, c clock.Clock, req *core.Record) (warnings []string, err error)
	// OnFinalApproved runs when the request is approved as final (leave: writes ledger debits).
	OnFinalApproved(txApp core.App, req *core.Record) error
	// OnCancelledAfterApproval runs when an approved request is cancelled (leave: writes reversals).
	OnCancelledAfterApproval(txApp core.App, req *core.Record) error
	// Summary is one line describing the request for inboxes and notifications.
	Summary(app core.App, req *core.Record) (string, error)
}

// UserError is a Validate error meant for the user; its message is shown as-is.
type UserError struct{ Message string }

func (e *UserError) Error() string { return e.Message }

const formsKey = "approvals.forms"

// RegisterForm makes f the hooks for requests of formType on app. Call it from the form module's Register.
func RegisterForm(app core.App, formType string, f Form) {
	forms, _ := app.Store().Get(formsKey).(map[string]Form)
	if forms == nil {
		forms = map[string]Form{}
		app.Store().Set(formsKey, forms)
	}
	forms[formType] = f
}

// FormFor returns the hooks registered for formType on app, and false when there are none.
func FormFor(app core.App, formType string) (Form, bool) {
	forms, _ := app.Store().Get(formsKey).(map[string]Form)
	f, ok := forms[formType]
	return f, ok
}
