package approvals_test

import (
	"errors"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	"officeapp/internal/core/approvals"
	"officeapp/internal/core/clock"
	"officeapp/internal/testapp"
)

var submitNow = clock.Fixed(time.Date(2026, 10, 1, 10, 0, 0, 0, clock.IST))

// fakeForm stands in for a form module: it returns the warnings and error it is given
// and remembers the request Validate saw.
type fakeForm struct {
	warnings  []string
	err       error
	validated *core.Record
}

func (f *fakeForm) Validate(_ core.App, _ clock.Clock, req *core.Record) ([]string, error) {
	f.validated = req
	return f.warnings, f.err
}
func (*fakeForm) OnFinalApproved(core.App, *core.Record) error          { return nil }
func (*fakeForm) OnCancelledAfterApproval(core.App, *core.Record) error { return nil }
func (*fakeForm) Summary(core.App, *core.Record) (string, error)        { return "fake", nil }

func setUser(t *testing.T, app core.App, code, field string, value any) {
	t.Helper()
	u, err := app.FindFirstRecordByData("users", "employee_code", code)
	if err != nil {
		t.Fatal(err)
	}
	u.Set(field, value)
	if err := app.Save(u); err != nil {
		t.Fatal(err)
	}
}

func submitInput(t *testing.T, app core.App, requester, approver string) approvals.SubmitInput {
	return approvals.SubmitInput{
		FormType:    "fake",
		FormVersion: 1,
		Requester:   userID(t, app, requester),
		Approver:    userID(t, app, approver),
		SaveForm:    func(core.App, *core.Record) error { return nil },
	}
}

func newApp(t *testing.T) (*tests.TestApp, *fakeForm) {
	app := testapp.New(t)
	f := &fakeForm{}
	approvals.RegisterForm(app, "fake", f)
	return app, f
}

func steps(t *testing.T, app core.App, req string) []*core.Record {
	t.Helper()
	s, err := app.FindRecordsByFilter("approval_steps", "request = {:r}", "seq", 0, 0, dbx.Params{"r": req})
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func count(t *testing.T, app core.App, collection string) int64 {
	t.Helper()
	n, err := app.CountRecords(collection)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func wantUserError(t *testing.T, err error, msg string) {
	t.Helper()
	var ue *approvals.UserError
	if !errors.As(err, &ue) {
		t.Fatalf("got error %v, want a UserError %q", err, msg)
	}
	if ue.Message != msg {
		t.Fatalf("got message %q, want %q", ue.Message, msg)
	}
}

func TestSubmitCreatesAPendingRequestWithASubmittedStep(t *testing.T) {
	app, f := newApp(t)
	f.warnings = []string{"Short notice."}
	in := submitInput(t, app, "E001", "E002")
	in.FormVersion = 3
	var saved *core.Record
	in.SaveForm = func(txApp core.App, req *core.Record) error {
		if !txApp.IsTransactional() {
			t.Error("SaveForm did not get the transaction app")
		}
		saved = req
		return nil
	}

	req, err := approvals.SubmitRequest(app, submitNow, in)
	if err != nil {
		t.Fatal(err)
	}
	if saved == nil || saved.Id != req.Id || f.validated == nil || f.validated.Id != req.Id {
		t.Fatalf("SaveForm and Validate must see the saved request; got %v, %v", saved, f.validated)
	}

	stored, err := app.FindRecordById("requests", req.Id)
	if err != nil {
		t.Fatal(err)
	}
	for field, want := range map[string]string{
		"form_type": "fake", "requester": in.Requester, "status": "pending", "current_approver": in.Approver,
		"final_approver": "", "form_version": "3", "recorded_by_hr": "false",
	} {
		if got := stored.GetString(field); got != want {
			t.Errorf("request %s = %v, want %v", field, got, want)
		}
	}
	if got := stored.GetDateTime("submitted_at").Time(); !got.Equal(submitNow.Now()) {
		t.Errorf("submitted_at = %v, want %v", got, submitNow.Now())
	}

	s := steps(t, app, req.Id)
	if len(s) != 1 {
		t.Fatalf("got %d steps, want 1", len(s))
	}
	var warnings []string
	if err := s[0].UnmarshalJSONField("warnings", &warnings); err != nil {
		t.Fatal(err)
	}
	if s[0].GetInt("seq") != 1 || s[0].GetString("action") != "submitted" || s[0].GetString("actor") != in.Requester ||
		s[0].GetString("to_user") != in.Approver || len(warnings) != 1 || warnings[0] != "Short notice." {
		t.Errorf("step = %v, warnings %v", s[0], warnings)
	}
}

func TestSubmitRefusals(t *testing.T) {
	cases := []struct {
		name, requester, approver, msg string
		prepare                        func(t *testing.T, app core.App)
	}{
		{"first approver is the requester", "E002", "E002", "Choose an approver other than the employee who made the request.", nil},
		{"first approver cannot approve", "E001", "E004", "Choose an active approver.", nil},
		{"first approver is inactive", "E001", "E002", "Choose an active approver.",
			func(t *testing.T, app core.App) { setUser(t, app, "E002", "active", false) }},
		{"requester is inactive", "E001", "E002", "Your account is not active.",
			func(t *testing.T, app core.App) { setUser(t, app, "E001", "active", false) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			app, f := newApp(t)
			if c.prepare != nil {
				c.prepare(t, app)
			}
			_, err := approvals.SubmitRequest(app, submitNow, submitInput(t, app, c.requester, c.approver))
			wantUserError(t, err, c.msg)
			if f.validated != nil || count(t, app, "requests") != 0 {
				t.Error("a refused submit reached the form or saved a request")
			}
		})
	}
}

func TestSubmitRefusesAnUnknownApprover(t *testing.T) {
	app, _ := newApp(t)
	in := submitInput(t, app, "E001", "E002")
	in.Approver = "nosuchuser0000"
	_, err := approvals.SubmitRequest(app, submitNow, in)
	wantUserError(t, err, "Choose an active approver.")
}

func TestSubmitRefusesAFormTypeWithNoForm(t *testing.T) {
	app, _ := newApp(t)
	in := submitInput(t, app, "E001", "E002")
	in.FormType = "expense"
	if _, err := approvals.SubmitRequest(app, submitNow, in); err == nil || count(t, app, "requests") != 0 {
		t.Fatalf("submit of an unregistered form type: %v", err)
	}
}

func TestSubmitValidateErrorLeavesNoRowsBehind(t *testing.T) {
	app, f := newApp(t)
	f.err = &approvals.UserError{Message: "CL allows at most 2 consecutive days."}
	in := submitInput(t, app, "E001", "E002")
	in.SaveForm = func(txApp core.App, _ *core.Record) error {
		setUser(t, txApp, "E001", "department", "Changed by the form")
		return nil
	}

	_, err := approvals.SubmitRequest(app, submitNow, in)
	wantUserError(t, err, "CL allows at most 2 consecutive days.")
	if count(t, app, "requests") != 0 || count(t, app, "approval_steps") != 0 {
		t.Error("rows were left behind")
	}
	u, _ := app.FindFirstRecordByData("users", "employee_code", "E001")
	if u.GetString("department") != "Operations" {
		t.Error("the form's own writes were not rolled back")
	}
}

func TestSubmitSaveFormErrorLeavesNoRowsBehind(t *testing.T) {
	app, f := newApp(t)
	in := submitInput(t, app, "E001", "E002")
	in.SaveForm = func(core.App, *core.Record) error { return errors.New("disk full") }
	if _, err := approvals.SubmitRequest(app, submitNow, in); err == nil {
		t.Fatal("a SaveForm error was ignored")
	}
	if f.validated != nil || count(t, app, "requests") != 0 {
		t.Error("the submit went on after SaveForm failed")
	}
}

func submitted(t *testing.T, app core.App) *core.Record {
	t.Helper()
	req, err := approvals.SubmitRequest(app, submitNow, submitInput(t, app, "E001", "E002"))
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func forwardInput(t *testing.T, app core.App, req *core.Record, actor, to string) approvals.ForwardInput {
	return approvals.ForwardInput{Request: req.Id, Actor: userID(t, app, actor), To: userID(t, app, to), Comment: "Over to you."}
}

func TestForwardMovesTheRequestToTheNextApprover(t *testing.T) {
	app, _ := newApp(t)
	req := submitted(t, app)
	in := forwardInput(t, app, req, "E002", "E003")
	if err := approvals.ForwardRequest(app, in); err != nil {
		t.Fatal(err)
	}

	stored, _ := app.FindRecordById("requests", req.Id)
	if stored.GetString("status") != "pending" || stored.GetString("current_approver") != in.To {
		t.Errorf("request = %v", stored)
	}
	s := steps(t, app, req.Id)
	if len(s) != 2 {
		t.Fatalf("got %d steps, want 2", len(s))
	}
	if s[1].GetInt("seq") != 2 || s[1].GetString("action") != "forwarded" || s[1].GetString("actor") != in.Actor ||
		s[1].GetString("to_user") != in.To || s[1].GetString("comment") != "Over to you." {
		t.Errorf("step = %v", s[1])
	}
}

func TestForwardRefusals(t *testing.T) {
	cases := []struct {
		name, actor, to, msg string
		prepare              func(t *testing.T, app core.App, req *core.Record)
	}{
		{"actor is not the current approver", "E003", "E002", "Only the current approver can forward this request.", nil},
		{"actor is the requester", "E001", "E003", "Only the current approver can forward this request.", nil},
		{"current approver is inactive", "E002", "E003", "Only the current approver can forward this request.",
			func(t *testing.T, app core.App, _ *core.Record) { setUser(t, app, "E002", "active", false) }},
		{"target is inactive", "E002", "E003", "Choose an active approver.",
			func(t *testing.T, app core.App, _ *core.Record) { setUser(t, app, "E003", "active", false) }},
		{"target cannot approve", "E002", "E004", "Choose an active approver.", nil},
		{"target is the actor", "E002", "E002", "Choose an approver other than yourself.", nil},
		{"target is the requester", "E002", "E001", "Choose an approver other than the employee who made the request.",
			func(t *testing.T, app core.App, _ *core.Record) { setUser(t, app, "E001", "can_approve", true) }},
		{"request is no longer pending", "E002", "E003", "This request is no longer pending.",
			func(t *testing.T, app core.App, req *core.Record) {
				req.Set("status", "approved")
				if err := app.Save(req); err != nil {
					t.Fatal(err)
				}
			}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			app, _ := newApp(t)
			req := submitted(t, app)
			if c.prepare != nil {
				c.prepare(t, app, req)
			}
			wantUserError(t, approvals.ForwardRequest(app, forwardInput(t, app, req, c.actor, c.to)), c.msg)
			if n := len(steps(t, app, req.Id)); n != 1 {
				t.Errorf("a refused forward recorded a step (%d steps)", n)
			}
			stored, _ := app.FindRecordById("requests", req.Id)
			if stored.GetString("current_approver") != userID(t, app, "E002") {
				t.Error("a refused forward changed the current approver")
			}
		})
	}
}

func TestForwardAllowsFiveForwardsAndRefusesTheSixth(t *testing.T) {
	app, _ := newApp(t)
	req := submitted(t, app)
	approvers := []string{"E002", "E003"}
	for i := range 5 {
		if err := approvals.ForwardRequest(app, forwardInput(t, app, req, approvers[i%2], approvers[(i+1)%2])); err != nil {
			t.Fatalf("forward %d: %v", i+1, err)
		}
	}
	s := steps(t, app, req.Id)
	if len(s) != 6 || s[5].GetInt("seq") != 6 {
		t.Fatalf("after 5 forwards got %d steps", len(s))
	}
	// After an odd number of forwards E003 holds the request.
	err := approvals.ForwardRequest(app, forwardInput(t, app, req, "E003", "E002"))
	wantUserError(t, err, "This request has already been forwarded 5 times.")
}
