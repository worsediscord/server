package util

import (
	"context"
	"log/slog"
)

type NopLogHandler struct{}

func (n NopLogHandler) Enabled(_ context.Context, _ slog.Level) bool  { return false }
func (n NopLogHandler) Handle(_ context.Context, _ slog.Record) error { return nil }
func (n NopLogHandler) WithAttrs(_ []slog.Attr) slog.Handler          { return n }
func (n NopLogHandler) WithGroup(_ string) slog.Handler               { return n }
