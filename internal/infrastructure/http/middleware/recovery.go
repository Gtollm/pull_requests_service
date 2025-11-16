package middleware

import (
	"fmt"
	"net/http"

	"pull-request-review/internal/infrastructure/adapters/logger"
)

func Recovery(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if err := recover(); err != nil {
						log.Error(
							fmt.Errorf("panic: %v", err),
							fmt.Sprintf("Panic recovered: %s %s", r.Method, r.URL.Path),
						)

						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						_, err := w.Write([]byte(`{"error":{"code":"INTERNAL_ERROR","message":"internal server error"}}`))
						if err != nil {
							return
						}
					}
				}()

				next.ServeHTTP(w, r)
			},
		)
	}
}