package logger

import (
	"log/slog"
	"os"
)

func New(envName string) *slog.Logger {
	var handler slog.Handler
	switch envName {
	case "dev":
		options := &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
		handler = slog.NewTextHandler(os.Stdout, options)
	default:
		panic("todo")
	}
	return slog.New(handler)
}
