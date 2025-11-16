package middleware

import (
	"fmt"
	"net/http"
	"time"

	"pull-request-review/internal/infrastructure/adapters/logger"
)

func Logger(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
				next.ServeHTTP(wrapped, r)
				duration := time.Since(start)
				log.Info(
					fmt.Sprintf("HTTP request: %s %s - %d (%s)", r.Method, r.URL.Path, wrapped.statusCode, duration),
				)
			},
		)
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}