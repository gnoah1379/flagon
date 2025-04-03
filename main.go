package main

import (
	"flagon/cmd"
	"log/slog"
)

// @title Flagon API
// @version 1.0
// @description API server for Flagon application
// @BasePath /api/v1
func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("error executing command", "error", err)
	}
}
