package app

import (
	"avito_pvz/internal/config"
	"avito_pvz/internal/repository"
	"avito_pvz/internal/service"
	"context"
	"log/slog"

	grpcapp "avito_pvz/internal/app/grpc"
	httpapp "avito_pvz/internal/app/http"

	httpserver "avito_pvz/internal/http"

	pgrepo "avito_pvz/internal/repository/pg"

	postgres "avito_pvz/internal/storage/pg"
)

type App struct {
	grpcServer *grpcapp.App
	httpServer *httpapp.App
}

func New(ctx context.Context, cfg config.Config, log *slog.Logger) *App {
	db := postgres.MustSetup(ctx, cfg.DB.DSN(), log)

	userRepo := repository.NewUser(pgrepo.NewPgUser(db))
	productRepo := repository.NewProduct(pgrepo.NewPgProduct(db))
	pvzRepo := repository.NewPVZ(pgrepo.NewPgPvz(db))
	receptionRepo := repository.NewReception(pgrepo.NewPgReception(db))

	productService := service.NewProduct(productRepo, receptionRepo, pvzRepo)
	pvzService := service.NewPVZServce(pvzRepo)
	receptionService := service.NewReceptionService(receptionRepo, pvzRepo)
	jwtService := service.NewJWTManager(cfg.JWT.SecretKey, cfg.JWT.Expire)
	userService := service.NewUserService(userRepo, jwtService)

	hndler := httpserver.NewServer(
		jwtService,
		userService,
		pvzService,
		receptionService,
		productService,
	)

	httpPvz := httpapp.NewApp(hndler, log)

	grpcPVZ := grpcapp.New(log, pvzService, cfg.GRPC.Port)

	return &App{
		grpcServer: grpcPVZ,
		httpServer: httpPvz,
	}
}

func (a App) Run() {
	go a.grpcServer.MustRun()
	a.httpServer.Run()
}
