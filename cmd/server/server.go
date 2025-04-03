package server

import (
	"context"
	"flagon/pkg/server"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Start the web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		srvCmd, err := New()
		if err != nil {
			return fmt.Errorf("failed to create server command: %w", err)
		}
		return srvCmd.Run()
	},
}

type CmdRunner struct {
	HttpServer *server.HttpServer
}

func (c *CmdRunner) Run() error {
	if err := c.HttpServer.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	signalChan := make(chan os.Signal, 1)
	defer signal.Stop(signalChan)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.HttpServer.Stop(ctx); err != nil {
		slog.Error("Failed to stop server", "error", err)
	}
	return nil
}
