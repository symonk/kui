package log

import (
	"log/slog"
)

// New instantiates a new (non global) logger instance.
func New(h slog.Handler) *slog.Logger {
	l := slog.New(h)
	return l
}
