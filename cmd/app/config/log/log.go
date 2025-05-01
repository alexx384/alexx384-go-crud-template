package log

import (
	slogctx "github.com/veqryn/slog-context"
	"log/slog"
	"os"
)

func CreateLogger() (*slog.Logger, *slog.LevelVar) {
	// Add a few default environmental attributes that always are included
	defaultAttrs := []slog.Attr{
		slog.String("service", "userService"),
	}
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}).WithAttrs(defaultAttrs)
	customHandler := slogctx.NewHandler(jsonHandler, nil)
	logger := slog.New(customHandler)
	slog.SetDefault(logger)
	return logger, logLevel
}
