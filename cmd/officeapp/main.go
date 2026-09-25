// Command officeapp runs the office forms and approvals app.
package main

import (
	"log"
	_ "time/tzdata" // Asia/Kolkata must resolve even on a machine without zoneinfo

	"github.com/pocketbase/pocketbase"

	"officeapp/internal/modules"
)

func main() {
	app := pocketbase.New()
	modules.Register(app, app.RootCmd)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
