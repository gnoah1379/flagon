package log

import (
	"flagon/pkg/config"
	"fmt"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"strings"
)

func Init() error {
	var cfg config.Log
	err := viper.Sub("log").Unmarshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to parse log config: %w", err)
	}
	var level slog.Level
	if err = level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}
	var logFile *os.File
	file := viper.GetString("log.file")
	switch strings.ToLower(file) {
	case "stderr":
		logFile = os.Stderr
	case "stdout":
		logFile = os.Stdout
	default:
		if logFile, err = os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
	}
	var opts = slog.HandlerOptions{
		AddSource: viper.GetBool("log.addSource"),
		Level:     level,
	}

	var handler slog.Handler
	format := strings.ToLower(viper.GetString("log.format"))
	switch format {
	case "json":
		handler = slog.NewJSONHandler(logFile, &opts)
	case "text":
		handler = slog.NewTextHandler(logFile, &opts)
	}
	slog.SetDefault(slog.New(handler))
	return nil
}
