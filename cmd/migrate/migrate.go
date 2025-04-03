package migrate

import (
	"errors"
	"flagon/pkg/migrations"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		migrateCmd, err := New()
		if err != nil {
			return fmt.Errorf("failed to create migrate command: %w", err)
		}
		return migrateCmd.Run()
	},
}

type CmdRunner struct {
	Migrator *migrations.Migrator
}

func (c *CmdRunner) Run() error {
	if err := c.Migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}
