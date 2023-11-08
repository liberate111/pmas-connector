package logger

import (
	"app-connector/config"
	"app-connector/constant"
	"fmt"

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
			handler.WithBackupNum(0),
			handler.WithBackupTime(0),
		)
	} else {
		h = handler.MustFileHandler("./log/pmas-connector-info.log",
			handler.WithLogLevels(slog.Levels{slog.InfoLevel, slog.ErrorLevel, slog.FatalLevel}),
			handler.WithMaxSize(10485760),
			handler.WithRotateTime(2629746),
			handler.WithBackupNum(0),
			handler.WithBackupTime(0),
		)
	}

	slog.PushHandler(h)

	Info("event: initial_logger, status: success")
	Debug("Initial config", fmt.Sprintf("config: %+v", config.Config))
	Info(constant.BANNER)
	Info("event: start_app, status: process, msg: start PMAS-CONNECTOR application...")
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
