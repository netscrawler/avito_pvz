package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	requestIDKey = "request_id"
	loggerKey    = "logger"
)

// LoggingMiddleware creates a structured logger and adds it to the context.
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := uuid.New().String()

			// Create a child logger with request context
			childLogger := logger.With(
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
			)

			// Add logger and request ID to context
			ctx := context.WithValue(r.Context(), loggerKey, childLogger)
			ctx = context.WithValue(ctx, requestIDKey, requestID)
			r = r.WithContext(ctx)

			// Log request start
			childLogger.Info("request started")

			// Create a custom response writer to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Call next handler
			next.ServeHTTP(rw, r)

			// Log request completion
			childLogger.Info("request completed",
				slog.Int("status_code", rw.statusCode),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

// AuthMiddleware checks for valid JWT token in Authorization header.
func AuthMiddleware(exceptPaths map[string]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if exceptPaths[r.URL.Path] {
				next.ServeHTTP(w, r)

				return
			}

			logger, ok := r.Context().Value(loggerKey).(*slog.Logger)
			if !ok {
				logger = slog.Default()
			}

			token := r.Header.Get("Authorization")
			if token == "" {
				logger.Error("missing authorization token")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)

				return
			}

			// TODO: Validate JWT token here
			// For now, we'll just pass it through

			next.ServeHTTP(w, r)
		})
	}
}

// TracingMiddleware adds tracing context to the request.
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get or create request ID
		requestID, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			requestID = uuid.New().String()
		}

		// Add tracing headers
		w.Header().Set("X-Request-ID", requestID)
		w.Header().Set("X-Trace-ID", requestID)

		next.ServeHTTP(w, r)
	})
}

// responseWriter is a custom response writer to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
