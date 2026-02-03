package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging returns a middleware that logs HTTP requests.
func Logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			msg := fmt.Sprintf("%s %s %d %s",
				r.Method, r.URL.Path, wrapped.statusCode, duration.String())
			logger.Info(r.Context(), msg)
		})
	}
}
