package approvals_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	"officeapp/internal/core/approvals"
	"officeapp/internal/testapp"
)

func userID(t *testing.T, app core.App, code string) string {
	t.Helper()
	u, err := app.FindFirstRecordByData("users", "employee_code", code)
	if err != nil {
		t.Fatal(err)
	}
	return u.Id
}

func newRequest(t *testing.T, app *tests.TestApp) *core.Record {
	t.Helper()
	c, err := app.FindCollectionByNameOrId("requests")
	if err != nil {
		t.Fatal(err)
	}
	r := core.NewRecord(c)
	r.Set("form_type", "leave")
	r.Set("requester", userID(t, app, "E001"))
	r.Set("status", string(approvals.Pending))
	r.Set("current_approver", userID(t, app, "E002"))
	r.Set("submitted_at", "2026-10-01 04:30:00.000Z")
	r.Set("form_version", 1)
	if err := app.Save(r); err != nil {
		t.Fatalf("saving a request: %v", err)
	}
	return r
}

func newStep(t *testing.T, app *tests.TestApp, req *core.Record, seq int, action approvals.StepAction) (*core.Record, error) {
	t.Helper()
	c, err := app.FindCollectionByNameOrId("approval_steps")
	if err != nil {
		t.Fatal(err)
	}
	s := core.NewRecord(c)
	s.Set("request", req.Id)
	s.Set("seq", seq)
	s.Set("actor", userID(t, app, "E001"))
	s.Set("action", string(action))
	s.Set("warnings", []string{"Short notice."})
	return s, app.Save(s)
}

func TestCollectionsAreClosedToTheAPI(t *testing.T) {
	app := testapp.New(t)
	for _, name := range []string{"requests", "approval_steps"} {
		c, err := app.FindCollectionByNameOrId(name)
		if err != nil {
			t.Fatal(err)
		}
		if c.ListRule != nil || c.ViewRule != nil || c.CreateRule != nil || c.UpdateRule != nil || c.DeleteRule != nil {
			t.Errorf("%s API rules must all be nil (superuser only)", name)
		}
	}
	requests, _ := app.FindCollectionByNameOrId("requests")
	for _, f := range []string{"form_type", "requester", "status", "current_approver", "final_approver", "submitted_at", "form_version", "recorded_by_hr"} {
		if requests.Fields.GetByName(f) == nil {
			t.Errorf("requests has no field %q", f)
		}
	}
}

func TestRequestStatusMustBeADesignStatus(t *testing.T) {
	app := testapp.New(t)
	r := newRequest(t, app)
	r.Set("status", "archived")
	if err := app.Save(r); err == nil {
		t.Error("a request was saved with status \"archived\"")
	}
	for _, s := range []approvals.Status{approvals.Approved, approvals.Rejected, approvals.Cancelled, approvals.CancelRequested} {
		r.Set("status", string(s))
		if err := app.Save(r); err != nil {
			t.Errorf("status %q refused: %v", s, err)
		}
	}
}

func TestApprovalStepsAreAppendOnly(t *testing.T) {
	app := testapp.New(t)
	req := newRequest(t, app)

	step, err := newStep(t, app, req, 1, approvals.StepSubmitted)
	if err != nil {
		t.Fatalf("saving a new step: %v", err)
	}
	if _, err := newStep(t, app, req, 2, approvals.StepForwarded); err != nil {
		t.Fatalf("saving the next step: %v", err)
	}
	if _, err := newStep(t, app, req, 1, approvals.StepRejected); err == nil {
		t.Error("a second step with seq 1 on the same request was saved")
	}
	if _, err := newStep(t, app, newRequest(t, app), 1, approvals.StepSubmitted); err != nil {
		t.Errorf("seq 1 on another request refused: %v", err)
	}
	if _, err := newStep(t, app, req, 3, approvals.StepAction("edited")); err == nil {
		t.Error("a step with an unknown action was saved")
	}

	step.Set("comment", "changed later")
	if err := app.Save(step); err == nil {
		t.Error("an approval step was updated")
	}
	if err := app.Delete(step); err == nil {
		t.Error("an approval step was deleted")
	}
	stored, err := app.FindRecordById("approval_steps", step.Id)
	if err != nil || stored.GetString("comment") != "" || stored.GetString("action") != string(approvals.StepSubmitted) {
		t.Errorf("stored step changed: %v, %v", stored, err)
	}
}

func TestStepActionValuesMatchTheDesign(t *testing.T) {
	got := []approvals.StepAction{approvals.StepSubmitted, approvals.StepForwarded, approvals.StepApprovedFinal, approvals.StepRejected,
		approvals.StepCancelled, approvals.StepCancelRequested, approvals.StepCancelApproved, approvals.StepCancelDeclined,
		approvals.StepReassigned, approvals.StepRecordedByHR}
	want := []string{"submitted", "forwarded", "approved_final", "rejected", "cancelled", "cancel_requested",
		"cancel_approved", "cancel_declined", "reassigned", "recorded_by_hr"}
	for i := range want {
		if string(got[i]) != want[i] {
			t.Errorf("step action %d = %q, want %q", i, got[i], want[i])
		}
	}
	app := testapp.New(t)
	req := newRequest(t, app)
	for i, a := range got {
		if _, err := newStep(t, app, req, i+1, a); err != nil {
			t.Errorf("step action %q refused by the collection: %v", a, err)
		}
	}
}
