package cmd

import (
	"context"
	"errors"
	"flagon/pkg/config"
	"flagon/server"
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, ok := cmd.Context().Value("config").(config.Config)
		if !ok {
			return fmt.Errorf("config not found")
		}
		srv, err := server.NewHttpServer(cfg.Server)
		if err != nil {
			return fmt.Errorf("failed to create server: %w", err)
		}
		go func() {
			if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("Failed to start server:", err)
			}
		}()
		signalChan := make(chan os.Signal, 1)
		defer signal.Stop(signalChan)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
		slog.Info("Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = srv.Stop(ctx); err != nil {
			slog.Error("Failed to stop server:", err)
		}
		return nil
	},
}
