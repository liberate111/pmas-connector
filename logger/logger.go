package logger

import (
	"app-connector/config"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

func InitLog() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
	})

	var h *handler.SyncCloseHandler
	if config.Config.Log.Level == "debug" {
		h = handler.MustFileHandler("./log/pmas-connector-debug.log",
			handler.WithLogLevels(slog.AllLevels),
			handler.WithMaxSize(10485760),
			handler.WithRotateTime(2629746),
		)
	} else {
		h = handler.MustFileHandler("./log/pmas-connector-info.log",
			handler.WithLogLevels(slog.Levels{slog.InfoLevel, slog.ErrorLevel, slog.FatalLevel}),
			handler.WithMaxSize(10485760),
			handler.WithRotateTime(2629746),
		)
	}

	slog.PushHandler(h)

	Info("Initial logger")
	Debug("Initial config", "config", config.Config)
	Info("Start PMAS-CONNECTOR Application...")
}

func Info(args ...any) {
	slog.Info(args)
}

func Error(args ...any) {
	slog.Error(args)
}

func Warn(args ...any) {
	slog.Warn(args)
}

func Fatal(args ...any) {
	slog.Fatal(args)
}

func Debug(args ...any) {
	slog.Debug(args)
}
