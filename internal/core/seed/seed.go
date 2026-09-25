// Package seed fills an empty database with demo data for development and previews.
// Each module adds its own seeder in its Register; `officeapp seed` runs them in that order.
package seed

import (
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

// Func writes one module's demo data inside the seed transaction.
type Func func(txApp core.App) error

type seeder struct {
	name string
	fn   Func
}

const storeKey = "seed.seeders"

// Add registers a module's seeder on app. Call it from the module's Register.
func Add(app core.App, name string, fn Func) {
	list, _ := app.Store().Get(storeKey).([]seeder)
	app.Store().Set(storeKey, append(list, seeder{name, fn}))
}

// Run applies every seeder in one transaction. It refuses a database that
// already has users, so it can never plant demo passwords in real data.
func Run(app core.App) error {
	n, err := app.CountRecords("users")
	if err != nil {
		return err
	}
	if n > 0 {
		return errors.New("seed: database already has users; seed runs only on an empty database")
	}
	return app.RunInTransaction(func(txApp core.App) error {
		list, _ := app.Store().Get(storeKey).([]seeder)
		for _, s := range list {
			if err := s.fn(txApp); err != nil {
				return fmt.Errorf("seed %s: %w", s.name, err)
			}
		}
		return nil
	})
}

// Command is `officeapp seed`.
func Command(app core.App) *cobra.Command {
	return &cobra.Command{
		Use:   "seed",
		Short: "Fill an empty database with demo data",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.Bootstrap(); err != nil {
				return err
			}
			if err := app.RunAllMigrations(); err != nil {
				return err
			}
			if err := Run(app); err != nil {
				return err
			}
			cmd.Println("seeded demo data")
			return nil
		},
	}
}
