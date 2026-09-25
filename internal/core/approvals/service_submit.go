package approvals

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/clock"
)

// maxForwards is how many times one request may be forwarded (design spec §7 step 2).
const maxForwards = 5

// SubmitInput is a requester sending a new request of a form to a first approver.
type SubmitInput struct {
	FormType    string
	FormVersion int
	Requester   string // users id; the requester is the actor
	Approver    string // users id of the first approver
	// SaveForm saves the form's own data for the just-saved request, through txApp.
	SaveForm func(txApp core.App, req *core.Record) error
}

// SubmitRequest creates a pending request, lets the form save its data and validate it, and records
// the submitted step, all in one transaction: a refusal or a Validate error leaves nothing behind.
func SubmitRequest(app core.App, c clock.Clock, in SubmitInput) (*core.Record, error) {
	form, ok := FormFor(app, in.FormType)
	if !ok {
		return nil, fmt.Errorf("approvals: no form registered for %q", in.FormType)
	}
	var req *core.Record
	err := app.RunInTransaction(func(txApp core.App) error {
		requester, err := txApp.FindRecordById("users", in.Requester)
		if err != nil {
			return err
		}
		if !requester.GetBool("active") {
			return &UserError{"Your account is not active."}
		}
		to, err := next(None, Submit)
		if err != nil {
			return err
		}
		if err := checkApprover(txApp, in.Approver, in.Requester, in.Requester); err != nil {
			return err
		}

		requests, err := txApp.FindCollectionByNameOrId("requests")
		if err != nil {
			return err
		}
		req = core.NewRecord(requests)
		req.Set("form_type", in.FormType)
		req.Set("form_version", in.FormVersion)
		req.Set("requester", in.Requester)
		req.Set("status", string(to))
		req.Set("current_approver", in.Approver)
		req.Set("submitted_at", c.Now())
		if err := txApp.Save(req); err != nil {
			return err
		}
		if err := in.SaveForm(txApp, req); err != nil {
			return err
		}
		warnings, err := form.Validate(txApp, c, req)
		if err != nil {
			return err
		}
		return addStep(txApp, req, in.Requester, StepSubmitted, in.Approver, "", warnings)
	})
	if err != nil {
		return nil, err
	}
	return req, nil
}

// ForwardInput is the current approver approving a request and sending it on to another approver.
type ForwardInput struct {
	Request string // requests id
	Actor   string // users id of the approver forwarding it
	To      string // users id of the next approver
	Comment string
}

// ForwardRequest makes To the current approver of a pending request and records the forwarded step.
func ForwardRequest(app core.App, in ForwardInput) error {
	return app.RunInTransaction(func(txApp core.App) error {
		req, err := requestForCurrentApprover(txApp, in.Request, in.Actor, "forward")
		if err != nil {
			return err
		}
		to, err := next(Status(req.GetString("status")), Forward)
		if err != nil {
			return err
		}
		if err := checkApprover(txApp, in.To, req.GetString("requester"), in.Actor); err != nil {
			return err
		}
		forwards, err := txApp.CountRecords("approval_steps", dbx.HashExp{"request": req.Id, "action": string(StepForwarded)})
		if err != nil {
			return err
		}
		if forwards >= maxForwards {
			return &UserError{fmt.Sprintf("This request has already been forwarded %d times.", maxForwards)}
		}

		req.Set("status", string(to))
		req.Set("current_approver", in.To)
		if err := txApp.Save(req); err != nil {
			return err
		}
		return addStep(txApp, req, in.Actor, StepForwarded, in.To, in.Comment, nil)
	})
}

// next is Next for a service action: a move the table does not allow is a UserError.
func next(from Status, a Action) (Status, error) {
	to, ok := Next(from, a)
	if !ok {
		return "", &UserError{"This request is no longer pending."}
	}
	return to, nil
}

// checkApprover refuses sending a request to anyone but an active approver who is
// neither the employee who made it nor the actor.
func checkApprover(txApp core.App, id, requester, actor string) error {
	if id == requester {
		return &UserError{"Choose an approver other than the employee who made the request."}
	}
	if id == actor {
		return &UserError{"Choose an approver other than yourself."}
	}
	u, err := txApp.FindRecordById("users", id)
	if errors.Is(err, sql.ErrNoRows) {
		return &UserError{"Choose an active approver."}
	}
	if err != nil {
		return err
	}
	if !u.GetBool("active") || !u.GetBool("can_approve") {
		return &UserError{"Choose an active approver."}
	}
	return nil
}

// requestForCurrentApprover reads a request and refuses it unless actor is its active current approver;
// verb names the refused action in the message.
func requestForCurrentApprover(txApp core.App, requestID, actorID, verb string) (*core.Record, error) {
	req, err := txApp.FindRecordById("requests", requestID)
	if err != nil {
		return nil, err
	}
	actor, err := txApp.FindRecordById("users", actorID)
	if err != nil {
		return nil, err
	}
	if req.GetString("current_approver") != actor.Id || !actor.GetBool("active") {
		return nil, &UserError{"Only the current approver can " + verb + " this request."}
	}
	return req, nil
}

// addStep records the request's next approval step, numbered after the ones it already has.
func addStep(txApp core.App, req *core.Record, actor string, action StepAction, toUser, comment string, warnings []string) error {
	seq, err := txApp.CountRecords("approval_steps", dbx.HashExp{"request": req.Id})
	if err != nil {
		return err
	}
	steps, err := txApp.FindCollectionByNameOrId("approval_steps")
	if err != nil {
		return err
	}
	s := core.NewRecord(steps)
	s.Set("request", req.Id)
	s.Set("seq", seq+1)
	s.Set("actor", actor)
	s.Set("action", string(action))
	s.Set("to_user", toUser)
	s.Set("comment", comment)
	s.Set("warnings", warnings)
	return txApp.Save(s)
}
