package main

//nolint:gofumpt
import (
	"context"
	"log/slog"
	"os"

	"avito_pvz/internal/app"
	"avito_pvz/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.ENV)
	log.InfoContext(context.Background(), "HELLO WORLD")
	log.InfoContext(context.Background(), "Starting with config", "cfg", cfg)

	app := app.New(context.Background(), *cfg, log)
	app.Run()
}

//nolint:exhaustruct
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
