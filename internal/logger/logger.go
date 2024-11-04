package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	cfg "github.com/andreychh/snippetbox/internal/config"
)

type Logger struct {
	*slog.Logger
}

func New(config cfg.Logger) (*Logger, error) {
	options := &slog.HandlerOptions{
		AddSource: config.AddSource,
		Level:     config.Level,
	}

	var writer io.Writer
	switch config.Writer {
	case "stdout":
		writer = os.Stdout
	default:
		return nil, fmt.Errorf("unknown writer type %q", config.Writer)
	}

	var handler slog.Handler
	switch config.Handler {
	case "text":
		handler = slog.NewTextHandler(writer, options)
	case "json":
		handler = slog.NewJSONHandler(writer, options)
	default:
		return nil, fmt.Errorf("unknown handler type %q", config.Handler)
	}

	return &Logger{Logger: slog.New(handler)}, nil
}

func Default() *Logger {
	return &Logger{Logger: slog.Default()}
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}
