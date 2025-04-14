package httpapp

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"avito_pvz/internal/http/gen"
)

// App структура, которая содержит http сервер
type App struct {
	httpServer *http.Server
}

// NewApp создает экземпляр App с зависимостями и handler'ом
func NewApp(handler gen.StrictServerInterface) *App {
	// Swagger schema (для валидации запросов и регистрации роутов)
	swagger, err := gen.GetSwagger()
	if err != nil {
		log.Fatalf("failed to load swagger spec: %v", err)
	}
	swagger.Servers = nil // не проверяем server.url в схеме

	// Генерация handler'а
	openapiHandler := gen.NewStrictHandler(handler, nil)

	// Роутинг через стандартную библиотеку
	mux := http.NewServeMux()
	gen.HandlerFromMux(openapiHandler, mux)

	// Конфигурация http.Server
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return &App{
		httpServer: srv,
	}
}

// Run запускает сервер и включает graceful shutdown
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

	// Ожидаем сигнала для остановки
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
