package app

import (
	grpcapp "avito_pvz/internal/app/grpc"
	httpapp "avito_pvz/internal/app/http"
)

type App struct {
	grpcServer *grpcapp.App
	httpServer *httpapp.App
}
