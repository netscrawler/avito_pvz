package httpapp

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpserver "avito_pvz/internal/http"
	"avito_pvz/internal/http/gen"
)

// App структура, которая содержит http сервер
type App struct {
	httpServer *http.Server
}

// NewApp создает экземпляр App с зависимостями и handler'ом
func NewApp(handler gen.StrictServerInterface, log *slog.Logger) *App {
	// Swagger schema (для валидации запросов и регистрации роутов)
	swagger, err := gen.GetSwagger()
	if err != nil {
		panic("err")
	}

	swagger.Servers = nil
	openapiHandler := gen.NewStrictHandler(handler, nil)

	exceptPaths := map[string]bool{
		"/register":   true,
		"/login":      true,
		"/dummyLogin": true,
	}

	middlewareChain := httpserver.LoggingMiddleware(log)(
		httpserver.AuthMiddleware(exceptPaths)(
			httpserver.TracingMiddleware(
				gen.HandlerFromMux(openapiHandler, http.NewServeMux()),
			),
		),
	)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           middlewareChain,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return &App{
		httpServer: srv,
	}
}

func (app *App) Run() {
	// Канал для остановки
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервер в горутине
	go func() {
		log.Printf("Starting server on %s", app.httpServer.Addr)
		if err := app.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stopCh
	log.Println("Received shutdown signal, gracefully shutting down...")

	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Останавливаем сервер
	if err := app.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}
