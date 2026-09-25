package leave_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"

	"officeapp/internal/forms/leave"
	"officeapp/internal/testapp"
)

// approvalRequest saves a pending leave row in requests, the parent every leave request needs.
func approvalRequest(t *testing.T, app *tests.TestApp) *core.Record {
	t.Helper()
	emp, err := app.FindFirstRecordByData("users", "employee_code", "E001")
	if err != nil {
		t.Fatal(err)
	}
	col, err := app.FindCollectionByNameOrId("requests")
	if err != nil {
		t.Fatal(err)
	}
	r := core.NewRecord(col)
	r.Set("form_type", "leave")
	r.Set("requester", emp.Id)
	r.Set("status", "pending")
	r.Set("submitted_at", "2026-10-01 04:30:00.000Z")
	r.Set("form_version", 1)
	if err := app.Save(r); err != nil {
		t.Fatalf("saving a request: %v", err)
	}
	return r
}

func leaveRequest(t *testing.T, app *tests.TestApp, requestID string) *core.Record {
	t.Helper()
	cl, err := leave.RuleFor(app, "CL", "2026-10-12")
	if err != nil {
		t.Fatal(err)
	}
	col, err := app.FindCollectionByNameOrId("leave_requests")
	if err != nil {
		t.Fatal(err)
	}
	r := core.NewRecord(col)
	r.Set("request", requestID)
	r.Set("leave_type", cl.GetString("leave_type"))
	r.Set("from_date", "2026-10-12")
	r.Set("from_session", "full")
	r.Set("to_date", "2026-10-12")
	r.Set("to_session", "full")
	r.Set("days", 1)
	r.Set("leave_year", "2026-27")
	r.Set("rule", cl.Id)
	return r
}

func TestLeaveRequestLinkRequired(t *testing.T) {
	app := testapp.New(t)
	if err := app.Save(leaveRequest(t, app, "")); err == nil {
		t.Error("a leave request was saved without an approval request")
	}
	if err := app.Save(leaveRequest(t, app, "missingrequest1")); err == nil {
		t.Error("a leave request was saved pointing at a request that does not exist")
	}
	if err := app.Save(leaveRequest(t, app, approvalRequest(t, app).Id)); err != nil {
		t.Errorf("a leave request with its approval request: %v", err)
	}
}

func TestLeaveRequestLinkOnePerRequest(t *testing.T) {
	app := testapp.New(t)
	req := approvalRequest(t, app)
	if err := app.Save(leaveRequest(t, app, req.Id)); err != nil {
		t.Fatal(err)
	}
	if err := app.Save(leaveRequest(t, app, req.Id)); err == nil {
		t.Error("a second leave request was saved for the same approval request")
	}
	if err := app.Save(leaveRequest(t, app, approvalRequest(t, app).Id)); err != nil {
		t.Errorf("a leave request for another approval request: %v", err)
	}
}

// Cascade delete is off: removing a request never silently removes its leave request.
func TestLeaveRequestLinkNoCascadeDelete(t *testing.T) {
	app := testapp.New(t)
	req := approvalRequest(t, app)
	lr := leaveRequest(t, app, req.Id)
	if err := app.Save(lr); err != nil {
		t.Fatal(err)
	}
	if err := app.Delete(req); err == nil {
		t.Error("a request with a leave request was deleted")
	}
	if _, err := app.FindRecordById("leave_requests", lr.Id); err != nil {
		t.Errorf("the leave request is gone: %v", err)
	}
}

func TestLeaveRequestLinkLedgerOptional(t *testing.T) {
	app := testapp.New(t)
	without := ledgerEntry(t, app, "2026-27", "adjustment", "", 1)
	if err := app.Save(without); err != nil {
		t.Errorf("a ledger entry without a request: %v", err)
	}
	req := approvalRequest(t, app)
	with := ledgerEntry(t, app, "2026-27", "debit", "", -1)
	with.Set("request", req.Id)
	if err := app.Save(with); err != nil {
		t.Fatalf("a ledger entry with a request: %v", err)
	}
	stored, err := app.FindRecordById("leave_ledger", with.Id)
	if err != nil {
		t.Fatal(err)
	}
	if stored.GetString("request") != req.Id {
		t.Errorf("stored request = %q, want %q", stored.GetString("request"), req.Id)
	}
	bad := ledgerEntry(t, app, "2026-27", "debit", "", -1)
	bad.Set("request", "missingrequest1")
	if err := app.Save(bad); err == nil {
		t.Error("a ledger entry was saved pointing at a request that does not exist")
	}
}
