package leave

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Links leave rows to their approval request (design §5.2): every leave request belongs to exactly one
// request; a ledger entry may name the request it came from. Cascade delete stays off, so a request
// with leave rows cannot be deleted.
func init() {
	m.Register(func(app core.App) error {
		requests, err := app.FindCollectionByNameOrId("requests")
		if err != nil {
			return err
		}

		leaveRequests, err := app.FindCollectionByNameOrId("leave_requests")
		if err != nil {
			return err
		}
		leaveRequests.Fields.Add(&core.RelationField{Name: "request", CollectionId: requests.Id, MaxSelect: 1, Required: true})
		leaveRequests.AddIndex("idx_leave_requests_request", true, "request", "")
		if err := app.Save(leaveRequests); err != nil {
			return err
		}

		ledger, err := app.FindCollectionByNameOrId("leave_ledger")
		if err != nil {
			return err
		}
		ledger.Fields.Add(&core.RelationField{Name: "request", CollectionId: requests.Id, MaxSelect: 1})
		return app.Save(ledger)
	}, nil)
}
