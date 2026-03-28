package core

import (
	"io"
	"log/slog"
)

// InitLogging configures the global slog.Default() logger with a JSON handler
// at the specified level writing to the specified output.
// It should be called once at application startup, before any goroutines log.
func InitLogging(level slog.Level, output io.Writer) {
	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: level,
	})
	slog.SetDefault(slog.New(handler))
}
