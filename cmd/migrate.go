package cmd

import (
	"flagon/internal/migrations"
	"flagon/pkg/config"
	"flagon/pkg/database"
	"fmt"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, ok := cmd.Context().Value("config").(config.Config)
		if !ok {
			return fmt.Errorf("config not found")
		}
		db, err := database.Open(cfg.Database)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer db.Close()

		return migrations.Up(cfg.Database.Driver, db)
	},
}
