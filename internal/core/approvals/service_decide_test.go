package approvals_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	"officeapp/internal/core/approvals"
)

// hookForm is a fakeForm that counts OnFinalApproved calls and can make them fail.
type hookForm struct {
	fakeForm
	approved    int
	approvedErr error
}

func (f *hookForm) OnFinalApproved(txApp core.App, req *core.Record) error {
	if !txApp.IsTransactional() {
		return errors.New("OnFinalApproved did not get the transaction app")
	}
	f.approved++
	return f.approvedErr
}

// pendingAtE002 is a request E001 submitted to E002, with the form's call records reset.
func pendingAtE002(t *testing.T) (*tests.TestApp, *hookForm, *core.Record) {
	app, _ := newApp(t)
	f := &hookForm{}
	approvals.RegisterForm(app, "fake", f)
	req := submitted(t, app)
	f.validated = nil
	return app, f, req
}

func decideInput(t *testing.T, app core.App, req *core.Record, actor, comment string) approvals.DecideInput {
	return approvals.DecideInput{Request: req.Id, Actor: userID(t, app, actor), Comment: comment}
}

func stored(t *testing.T, app core.App, req *core.Record) *core.Record {
	t.Helper()
	r, err := app.FindRecordById("requests", req.Id)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestApproveFinalApprovesAndCallsTheHookOnce(t *testing.T) {
	app, f, req := pendingAtE002(t)
	in := decideInput(t, app, req, "E002", "Enjoy.")
	if err := approvals.ApproveFinalRequest(app, submitNow, in); err != nil {
		t.Fatal(err)
	}

	if f.validated == nil || f.validated.Id != req.Id {
		t.Error("approve as final did not re-run Validate")
	}
	if f.approved != 1 {
		t.Errorf("OnFinalApproved called %d times, want 1", f.approved)
	}
	r := stored(t, app, req)
	if r.GetString("status") != "approved" || r.GetString("final_approver") != in.Actor || r.GetString("current_approver") != "" {
		t.Errorf("request = %v", r)
	}
	s := steps(t, app, req.Id)
	if len(s) != 2 {
		t.Fatalf("got %d steps, want 2", len(s))
	}
	if s[1].GetInt("seq") != 2 || s[1].GetString("action") != "approved_final" || s[1].GetString("actor") != in.Actor ||
		s[1].GetString("comment") != "Enjoy." || s[1].GetString("to_user") != "" {
		t.Errorf("step = %v", s[1])
	}
}

func TestApproveFinalTwiceIsRefusedAndTheHookRunsOnce(t *testing.T) {
	app, f, req := pendingAtE002(t)
	in := decideInput(t, app, req, "E002", "")
	if err := approvals.ApproveFinalRequest(app, submitNow, in); err != nil {
		t.Fatal(err)
	}
	wantUserError(t, approvals.ApproveFinalRequest(app, submitNow, in), "Only the current approver can approve this request.")
	if f.approved != 1 || len(steps(t, app, req.Id)) != 2 {
		t.Errorf("second approve: hook count %d, %d steps", f.approved, len(steps(t, app, req.Id)))
	}
}

func TestApproveFinalValidateErrorLeavesTheRequestPending(t *testing.T) {
	app, f, req := pendingAtE002(t)
	f.err = &approvals.UserError{Message: "Not enough CL balance."}
	wantUserError(t, approvals.ApproveFinalRequest(app, submitNow, decideInput(t, app, req, "E002", "")), "Not enough CL balance.")
	assertUndecided(t, app, req)
	if f.approved != 0 {
		t.Error("OnFinalApproved ran after a Validate error")
	}
}

func TestApproveFinalHookErrorRollsBack(t *testing.T) {
	app, f, req := pendingAtE002(t)
	f.approvedErr = errors.New("ledger write failed")
	if err := approvals.ApproveFinalRequest(app, submitNow, decideInput(t, app, req, "E002", "")); err == nil {
		t.Fatal("an OnFinalApproved error was ignored")
	}
	assertUndecided(t, app, req)
}

func TestApproveFinalRefusesAFormTypeWithNoForm(t *testing.T) {
	app, _, req := pendingAtE002(t)
	req.Set("form_type", "expense")
	if err := app.Save(req); err != nil {
		t.Fatal(err)
	}
	if err := approvals.ApproveFinalRequest(app, submitNow, decideInput(t, app, req, "E002", "")); err == nil {
		t.Fatal("approved a request whose form type has no form")
	}
	assertUndecided(t, app, req)
}

func TestRejectRejectsWithTheComment(t *testing.T) {
	app, f, req := pendingAtE002(t)
	in := decideInput(t, app, req, "E002", "Team is short that week.")
	if err := approvals.RejectRequest(app, in); err != nil {
		t.Fatal(err)
	}
	if r := stored(t, app, req); r.GetString("status") != "rejected" || r.GetString("final_approver") != "" {
		t.Errorf("request = %v", r)
	}
	if f.approved != 0 || f.validated != nil {
		t.Error("reject ran the form's approval checks or hook")
	}
	s := steps(t, app, req.Id)
	if len(s) != 2 {
		t.Fatalf("got %d steps, want 2", len(s))
	}
	if s[1].GetInt("seq") != 2 || s[1].GetString("action") != "rejected" || s[1].GetString("actor") != in.Actor ||
		s[1].GetString("comment") != "Team is short that week." {
		t.Errorf("step = %v", s[1])
	}
}

// assertUndecided checks a refused action left the request pending at E002 with only its submitted step.
func assertUndecided(t *testing.T, app core.App, req *core.Record) {
	t.Helper()
	r := stored(t, app, req)
	if r.GetString("status") != "pending" || r.GetString("current_approver") != userID(t, app, "E002") || r.GetString("final_approver") != "" {
		t.Errorf("a refused action changed the request: %v", r)
	}
	if n := len(steps(t, app, req.Id)); n != 1 {
		t.Errorf("a refused action recorded a step (%d steps)", n)
	}
}

type refusal struct {
	name, actor, msg string
	prepare          func(t *testing.T, app core.App, req *core.Record)
}

// decideRefusals are refused for both approve as final and reject; VERB is the action in the message.
var decideRefusals = []refusal{
	{"actor is not the current approver", "E003", "Only the current approver can VERB this request.", nil},
	{"actor is the requester", "E001", "Only the current approver can VERB this request.", nil},
	{"current approver is inactive", "E002", "Only the current approver can VERB this request.",
		func(t *testing.T, app core.App, _ *core.Record) { setUser(t, app, "E002", "active", false) }},
	{"request is not pending", "E002", "This request is no longer pending.",
		func(t *testing.T, app core.App, req *core.Record) {
			req.Set("status", "cancel_requested")
			if err := app.Save(req); err != nil {
				t.Fatal(err)
			}
		}},
}

func TestApproveFinalRefusals(t *testing.T) {
	for _, c := range decideRefusals {
		t.Run(c.name, func(t *testing.T) {
			app, f, req := pendingAtE002(t)
			if c.prepare != nil {
				c.prepare(t, app, req)
			}
			err := approvals.ApproveFinalRequest(app, submitNow, decideInput(t, app, req, c.actor, ""))
			wantUserError(t, err, strings.ReplaceAll(c.msg, "VERB", "approve"))
			if f.approved != 0 || f.validated != nil || len(steps(t, app, req.Id)) != 1 {
				t.Error("a refused approve reached the form or recorded a step")
			}
			if r := stored(t, app, req); r.GetString("final_approver") != "" || r.GetString("current_approver") != userID(t, app, "E002") {
				t.Errorf("a refused approve changed the request: %v", r)
			}
		})
	}
}

func TestRejectRefusals(t *testing.T) {
	for _, c := range decideRefusals {
		t.Run(c.name, func(t *testing.T) {
			app, _, req := pendingAtE002(t)
			if c.prepare != nil {
				c.prepare(t, app, req)
			}
			wantUserError(t, approvals.RejectRequest(app, decideInput(t, app, req, c.actor, "No.")), strings.ReplaceAll(c.msg, "VERB", "reject"))
			assertNotRejected(t, app, req)
		})
	}
}

func TestRejectRefusesAnEmptyComment(t *testing.T) {
	app, _, req := pendingAtE002(t)
	wantUserError(t, approvals.RejectRequest(app, decideInput(t, app, req, "E002", "  ")), "Add a comment saying why the request is rejected.")
	assertUndecided(t, app, req)
}

func assertNotRejected(t *testing.T, app core.App, req *core.Record) {
	t.Helper()
	if len(steps(t, app, req.Id)) != 1 {
		t.Error("a refused reject recorded a step")
	}
	if stored(t, app, req).GetString("status") == "rejected" {
		t.Error("a refused reject rejected the request")
	}
}
