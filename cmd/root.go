package cmd

import (
	"flagon/cmd/migrate"
	"flagon/cmd/server"
	"flagon/pkg/config"
	"flagon/pkg/log"

	"github.com/spf13/cobra"
)

func Execute() error {
	return cmd.Execute()
}

var cfgFile string
var cmd = &cobra.Command{
	Use:   "flagon",
	Short: "Flagon is web application for managing feature flags for your projects.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := config.LoadConfig(cfgFile); err != nil {
			return err
		}
		if err := log.Init(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	cmd.AddCommand(server.Cmd)
	cmd.AddCommand(migrate.Cmd)

}
