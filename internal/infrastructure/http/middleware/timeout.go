package middleware

import (
	"context"
	"net/http"
	"time"
)

func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx, cancel := context.WithTimeout(r.Context(), timeout)
				defer cancel()

				r = r.WithContext(ctx)

				done := make(chan struct{})
				go func() {
					next.ServeHTTP(w, r)
					close(done)
				}()

				select {
				case <-done:
					return
				case <-ctx.Done():
					w.WriteHeader(http.StatusRequestTimeout)
					_, err := w.Write([]byte(`{"error":{"code":"TIMEOUT","message":"request timeout"}}`))
					if err != nil {
						return
					}
				}
			},
		)
	}
}