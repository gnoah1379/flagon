package cmd

import (
	"context"
	"flagon/pkg/config"
	"flagon/pkg/log"
	"github.com/spf13/cobra"
)

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "flagon",
	Short: "Flagon is web application for managing feature flags for your projects.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			return err
		}
		cmd.SetContext(context.WithValue(cmd.Context(), "config", cfg))
		if err := log.Setup(cfg.Log); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file")

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(migrateCmd)

}
