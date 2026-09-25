package approvals

import (
	"embed"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"

	"officeapp/internal/core/auth"
	"officeapp/internal/core/clock"
	"officeapp/internal/core/web"
)

//go:embed templates/my_requests.html
var myRequestsFS embed.FS

var myRequestsPages = web.NewPages(myRequestsFS)

func init() {
	addPart(func(app core.App) {
		app.OnServe().BindFunc(func(se *core.ServeEvent) error {
			se.Router.GET("/requests", showMyRequests).BindFunc(auth.RequireUser)
			return se.Next()
		})
	})
}

// statusLabels are the words the requester sees for each status.
var statusLabels = map[Status]string{
	Pending:         "Pending",
	Approved:        "Approved",
	Rejected:        "Rejected",
	Cancelled:       "Cancelled",
	CancelRequested: "Cancellation requested",
}

// myRequest is one row of the my requests page.
type myRequest struct {
	ID        string
	Summary   string
	Submitted string // date in IST, e.g. "2 Oct 2026"
	Status    string
	Approver  string // current approver's name; empty once the request is decided
}

func showMyRequests(e *core.RequestEvent) error {
	rows, err := myRequests(e.App, e.Auth.Id)
	if err != nil {
		return err
	}
	return myRequestsPages.Render(e, "my_requests", "My requests", rows)
}

// myRequests returns the requests of requester, newest first. It only reads.
func myRequests(app core.App, requester string) ([]myRequest, error) {
	records, err := app.FindRecordsByFilter("requests", "requester = {:u}", "-submitted_at", 0, 0, dbx.Params{"u": requester})
	if err != nil {
		return nil, err
	}
	if errs := app.ExpandRecords(records, []string{"current_approver"}, nil); len(errs) > 0 {
		return nil, fmt.Errorf("approvals: expanding current approvers: %v", errs)
	}
	rows := make([]myRequest, 0, len(records))
	for _, r := range records {
		form, ok := FormFor(app, r.GetString("form_type"))
		if !ok {
			return nil, fmt.Errorf("approvals: no form registered for %q", r.GetString("form_type"))
		}
		summary, err := form.Summary(app, r)
		if err != nil {
			return nil, err
		}
		row := myRequest{
			ID:        r.Id,
			Summary:   summary,
			Submitted: r.GetDateTime("submitted_at").Time().In(clock.IST).Format("2 Jan 2006"),
			Status:    statusLabels[Status(r.GetString("status"))],
		}
		if approver := r.ExpandedOne("current_approver"); approver != nil {
			row.Approver = approver.GetString("name")
		}
		rows = append(rows, row)
	}
	return rows, nil
}
