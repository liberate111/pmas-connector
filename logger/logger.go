package logger

import (
	"app-connector/config"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLog() {
	var opts *slog.HandlerOptions
	if config.Config.Log.Level == "debug" {
		opts = &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
	} else {
		opts = &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	Logger = slog.New(handler)
	Logger.Info("Initial logger")
	Logger.Info("Initial config", "config", config.Config)
	Logger.Info("Start PMAS-CONNECTOR Application...")
}
