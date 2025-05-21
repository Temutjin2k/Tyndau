package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// LoggingMiddleware wraps an http.Handler and logs requests
func LoggingMiddleware(next http.Handler, log *zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture the status code
		lrw := NewLoggingResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(lrw, r)

		// Log the request details
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", lrw.statusCode).
			Dur("duration", time.Since(start)).
			Str("remote_ip", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("request completed")
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture the status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK} // Default to 200 OK
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
