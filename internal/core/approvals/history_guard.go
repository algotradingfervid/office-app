package approvals

import (
	"errors"

	"github.com/pocketbase/pocketbase/core"
)

var errStepsAppendOnly = errors.New("approval steps are never edited or deleted; record a new step instead")

// Approval steps are the request's history (design spec §9 rule 1). Blocking the record hooks
// also stops edits from the superuser dashboard, not only from app code.
func init() {
	addPart(func(app core.App) {
		refuse := func(e *core.RecordEvent) error { return errStepsAppendOnly }
		app.OnRecordUpdate("approval_steps").BindFunc(refuse)
		app.OnRecordDelete("approval_steps").BindFunc(refuse)
	})
}
