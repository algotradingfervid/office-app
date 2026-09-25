package approvals_test

import (
	"testing"

	"officeapp/internal/core/approvals"
)

var (
	allStatuses = []approvals.Status{approvals.None, approvals.Pending, approvals.Approved, approvals.Rejected, approvals.Cancelled, approvals.CancelRequested}
	allActions  = []approvals.Action{approvals.Submit, approvals.RecordOnBehalf, approvals.Forward, approvals.ApproveFinal, approvals.Reject, approvals.Cancel,
		approvals.Reassign, approvals.RequestCancel, approvals.HRCancel, approvals.ApproveCancel, approvals.DeclineCancel}
)

// Every row of the design's §7 transition table ("pending / cancel_requested | reassign" is two rows).
var designTable = []struct {
	from approvals.Status
	a    approvals.Action
	to   approvals.Status
}{
	{approvals.None, approvals.Submit, approvals.Pending},
	{approvals.None, approvals.RecordOnBehalf, approvals.Approved},
	{approvals.Pending, approvals.Forward, approvals.Pending},
	{approvals.Pending, approvals.ApproveFinal, approvals.Approved},
	{approvals.Pending, approvals.Reject, approvals.Rejected},
	{approvals.Pending, approvals.Cancel, approvals.Cancelled},
	{approvals.Pending, approvals.Reassign, approvals.Pending},
	{approvals.CancelRequested, approvals.Reassign, approvals.CancelRequested},
	{approvals.Approved, approvals.RequestCancel, approvals.CancelRequested},
	{approvals.Approved, approvals.HRCancel, approvals.Cancelled},
	{approvals.CancelRequested, approvals.ApproveCancel, approvals.Cancelled},
	{approvals.CancelRequested, approvals.DeclineCancel, approvals.Approved},
}

func TestNextAllowsEveryRowOfTheTable(t *testing.T) {
	for _, r := range designTable {
		to, ok := approvals.Next(r.from, r.a)
		if !ok || to != r.to {
			t.Errorf("Next(%q, %q) = %q, %v; want %q, true", r.from, r.a, to, ok, r.to)
		}
	}
}

// Every (status, action) pair not in the table is refused, so no status can be reached another way.
func TestNextRefusesEverythingElse(t *testing.T) {
	inTable := map[[2]string]bool{}
	for _, r := range designTable {
		inTable[[2]string{string(r.from), string(r.a)}] = true
	}
	refused := 0
	for _, from := range allStatuses {
		for _, a := range allActions {
			if inTable[[2]string{string(from), string(a)}] {
				continue
			}
			refused++
			if to, ok := approvals.Next(from, a); ok || to != "" {
				t.Errorf("Next(%q, %q) = %q, %v; want refused", from, a, to, ok)
			}
		}
	}
	if want := len(allStatuses)*len(allActions) - len(designTable); refused != want {
		t.Fatalf("checked %d forbidden moves, want %d", refused, want)
	}
	if _, ok := approvals.Next(approvals.Pending, approvals.Action("delete")); ok {
		t.Error("an unknown action was allowed")
	}
}

func TestStatusValuesMatchTheDesign(t *testing.T) {
	want := []string{"", "pending", "approved", "rejected", "cancelled", "cancel_requested"}
	for i, s := range allStatuses {
		if string(s) != want[i] {
			t.Errorf("status %d = %q, want %q", i, s, want[i])
		}
	}
}
