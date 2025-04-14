package main

//nolint:gofumpt
import (
	"context"

	"avito_pvz/internal/app"
	"avito_pvz/internal/config"
	logger "avito_pvz/internal/pkg"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	logger.Init(cfg.ENV)
	log := logger.L()
	log.InfoContext(context.Background(), "HELLO WORLD")
	log.InfoContext(context.Background(), "Starting with config", "cfg", cfg)

	app := app.New(context.Background(), *cfg, log)
	app.Run()
}
