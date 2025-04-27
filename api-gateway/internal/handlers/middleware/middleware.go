package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Middleware func(next http.Handler) http.Handler

func NewMiddlewareChain(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

func NewTimeoutContextMW(timeoutInSec int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutInSec))
				defer cancel()

				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			})

	}
}

func NewLoggerMW(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				next.ServeHTTP(w, r)
				duration := time.Since(start)
				logger.Info("request",
					slog.String("method", r.Method),
					slog.String("url", r.URL.String()),
					slog.String("remote_addr", r.RemoteAddr),
					slog.String("user_agent", r.UserAgent()),
					slog.String("duration", duration.String()),
				)
			},
		)
	}
}

func NewCORS(CORS_URLS string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				allowedOrigins := strings.Split(CORS_URLS, ",")
				origin := r.Header.Get("Origin")

				allowOrigin := ""
				for _, o := range allowedOrigins {
					if strings.TrimSpace(o) == origin {
						allowOrigin = origin
						break
					}
				}

				if allowOrigin != "" {
					w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				// Handle preflight OPTIONS request
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusOK)
					return
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}

func RecoveryMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RequestValidator(msg proto.Message) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength == 0 {
				next.ServeHTTP(w, r)
				return
			}

			// Clone prototype to avoid mutation
			clone := proto.Clone(msg)

			decoder := protojson.UnmarshalOptions{
				DiscardUnknown: true,
			}

			var body []byte
			json.NewDecoder(r.Body).Decode(&body)
			if err := decoder.Unmarshal(body, clone); err != nil {
				http.Error(w, "Invalid request format", http.StatusBadRequest)
				return
			}

			// Add validation logic here using protoc-gen-validate
			if validator, ok := clone.(interface {
				Validate() error
			}); ok {
				if err := validator.Validate(); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

			// Store validated message in context
			ctx := context.WithValue(r.Context(), "validated_request", clone)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
