package approvals

import (
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/clock"
)

// DecideInput is the current approver approving a request as final or rejecting it.
type DecideInput struct {
	Request string // requests id
	Actor   string // users id of the approver deciding
	Comment string // required to reject
}

// ApproveFinalRequest re-runs the form's checks and, when they pass, approves a pending request,
// calls the form's OnFinalApproved and records the approved_final step, all in one transaction.
func ApproveFinalRequest(app core.App, c clock.Clock, in DecideInput) error {
	return app.RunInTransaction(func(txApp core.App) error {
		req, err := requestForCurrentApprover(txApp, in.Request, in.Actor, "approve")
		if err != nil {
			return err
		}
		to, err := next(Status(req.GetString("status")), ApproveFinal)
		if err != nil {
			return err
		}
		form, ok := FormFor(txApp, req.GetString("form_type"))
		if !ok {
			return fmt.Errorf("approvals: no form registered for %q", req.GetString("form_type"))
		}
		if _, err := form.Validate(txApp, c, req); err != nil {
			return err
		}

		req.Set("status", string(to))
		req.Set("final_approver", in.Actor)
		req.Set("current_approver", "")
		if err := txApp.Save(req); err != nil {
			return err
		}
		if err := form.OnFinalApproved(txApp, req); err != nil {
			return err
		}
		return addStep(txApp, req, in.Actor, StepApprovedFinal, "", in.Comment, nil)
	})
}

// RejectRequest rejects a pending request with the current approver's comment and records the rejected step.
func RejectRequest(app core.App, in DecideInput) error {
	if strings.TrimSpace(in.Comment) == "" {
		return &UserError{"Add a comment saying why the request is rejected."}
	}
	return app.RunInTransaction(func(txApp core.App) error {
		req, err := requestForCurrentApprover(txApp, in.Request, in.Actor, "reject")
		if err != nil {
			return err
		}
		to, err := next(Status(req.GetString("status")), Reject)
		if err != nil {
			return err
		}

		req.Set("status", string(to))
		if err := txApp.Save(req); err != nil {
			return err
		}
		return addStep(txApp, req, in.Actor, StepRejected, "", in.Comment, nil)
	})
}
