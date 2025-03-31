package main

import (
	"flagon/cmd"
	"log/slog"
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("error executing command", "error", err)
	}
}
