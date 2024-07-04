package util

import (
	"context"
	"log/slog"
	"strings"
)

var NopLogHandler nopLogHandler

type nopLogHandler struct{}

func (n nopLogHandler) Enabled(_ context.Context, _ slog.Level) bool  { return false }
func (n nopLogHandler) Handle(_ context.Context, _ slog.Record) error { return nil }
func (n nopLogHandler) WithAttrs(_ []slog.Attr) slog.Handler          { return n }
func (n nopLogHandler) WithGroup(_ string) slog.Handler               { return n }

// StringToLogLevel converts a string to a slog.Level. Defaults to slog.LevelInfo.
func StringToLogLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
