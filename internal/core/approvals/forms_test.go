package approvals_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/approvals"
	"officeapp/internal/core/clock"
	"officeapp/internal/testapp"
)

type stubForm struct{ name string }

func (stubForm) Validate(core.App, clock.Clock, *core.Record) ([]string, error) { return nil, nil }
func (stubForm) OnFinalApproved(core.App, *core.Record) error                   { return nil }
func (stubForm) OnCancelledAfterApproval(core.App, *core.Record) error          { return nil }
func (f stubForm) Summary(core.App, *core.Record) (string, error)               { return f.name, nil }

func TestRegisterFormRoundTrip(t *testing.T) {
	app := testapp.New(t)
	leave, other := stubForm{"leave"}, stubForm{"other"}
	approvals.RegisterForm(app, "leave", leave)
	approvals.RegisterForm(app, "other", other)

	for formType, want := range map[string]stubForm{"leave": leave, "other": other} {
		got, ok := approvals.FormFor(app, formType)
		if !ok || got != want {
			t.Errorf("FormFor(%q) = %v, %v; want %v, true", formType, got, ok, want)
		}
	}
	if f, ok := approvals.FormFor(app, "expense"); ok || f != nil {
		t.Errorf("FormFor(unregistered) = %v, %v; want nil, false", f, ok)
	}
	if _, ok := approvals.FormFor(testapp.New(t), "leave"); ok {
		t.Error("a form registered on one app is visible on another")
	}
}

func TestUserErrorShowsItsMessageAsIs(t *testing.T) {
	err := fmt.Errorf("validate: %w", &approvals.UserError{Message: "CL allows at most 2 consecutive days."})
	var ue *approvals.UserError
	if !errors.As(err, &ue) || ue.Error() != "CL allows at most 2 consecutive days." {
		t.Fatalf("got %v", ue)
	}
}
