package approvals

import (
	"embed"
	"fmt"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/auth"
	"officeapp/internal/core/clock"
	"officeapp/internal/core/web"
)

//go:embed templates/inbox.html
var inboxFS embed.FS

var inboxPages = web.NewPages(inboxFS)

func init() {
	addPart(func(app core.App) {
		app.OnServe().BindFunc(func(se *core.ServeEvent) error {
			se.Router.GET("/approvals", showInbox).BindFunc(auth.RequireUser)
			return se.Next()
		})
	})
}

// inboxRequest is one row of the pending my approval page.
type inboxRequest struct {
	ID        string
	Summary   string
	Requester string
	Submitted string // date in IST, e.g. "2 Oct 2026"
	Waiting   string // "today", "1 day", "3 days"
	Status    string
}

func showInbox(e *core.RequestEvent) error {
	if !e.Auth.GetBool("can_approve") {
		return e.ForbiddenError("Only approvers can see this page.", nil)
	}
	rows, err := inbox(e.App, clock.System{}, e.Auth.Id)
	if err != nil {
		return err
	}
	return inboxPages.Render(e, "inbox", "Pending my approval", rows)
}

// inbox returns the requests waiting for approver to act, oldest first. It only reads.
// Status is filtered too: a rejected request keeps its current_approver.
func inbox(app core.App, c clock.Clock, approver string) ([]inboxRequest, error) {
	records, err := app.FindRecordsByFilter("requests",
		"current_approver = {:u} && (status = {:pending} || status = {:cancel})", "submitted_at", 0, 0,
		dbx.Params{"u": approver, "pending": string(Pending), "cancel": string(CancelRequested)})
	if err != nil {
		return nil, err
	}
	if errs := app.ExpandRecords(records, []string{"requester"}, nil); len(errs) > 0 {
		return nil, fmt.Errorf("approvals: expanding requesters: %v", errs)
	}
	today := day(c.Now())
	rows := make([]inboxRequest, 0, len(records))
	for _, r := range records {
		form, ok := FormFor(app, r.GetString("form_type"))
		if !ok {
			return nil, fmt.Errorf("approvals: no form registered for %q", r.GetString("form_type"))
		}
		summary, err := form.Summary(app, r)
		if err != nil {
			return nil, err
		}
		submitted := r.GetDateTime("submitted_at").Time().In(clock.IST)
		rows = append(rows, inboxRequest{
			ID:        r.Id,
			Summary:   summary,
			Requester: r.ExpandedOne("requester").GetString("name"),
			Submitted: submitted.Format("2 Jan 2006"),
			Waiting:   waiting(int(today.Sub(day(submitted)).Hours() / 24)),
			Status:    statusLabels[Status(r.GetString("status"))],
		})
	}
	return rows, nil
}

// day is t's calendar date in IST at midnight UTC, so subtracting two days ignores clock changes.
func day(t time.Time) time.Time {
	y, m, d := t.In(clock.IST).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func waiting(days int) string {
	switch days {
	case 0:
		return "today"
	case 1:
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}
