package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey struct{}

var baseLogger *slog.Logger

func Init(env string) {
	switch env {
	case "local":
		baseLogger = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}),
		)
	case "prod":
		baseLogger = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		baseLogger = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelWarn}),
		)
	}
}

func L() *slog.Logger {
	return baseLogger
}

func WithCtx(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

func FromCtx(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(ctxKey{}).(*slog.Logger)
	if !ok {
		return baseLogger
	}
	return log
}
