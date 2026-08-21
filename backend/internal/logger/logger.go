package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey struct{}

func New(level slog.Level) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: level,
				ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						return slog.Attr{Key: "timestamp", Value: a.Value}
					}
					if a.Key == slog.MessageKey {
						return slog.Attr{Key: "message", Value: a.Value}
					}
					if a.Key == slog.LevelKey {
						return slog.Attr{Key: "leve", Value: a.Value}
					}
					return a
				},
			},
		),
	)
}

func WithContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok && log != nil {
		return log
	}
	return slog.Default()
}
