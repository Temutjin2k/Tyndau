package interseptor

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryLoggingInterceptor logs every gRPC request with zerolog
func UnaryLoggingInterceptor(log *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		// Handle the request
		resp, err = handler(ctx, req)

		duration := time.Since(start)

		// Initialize the log entry
		e := log.Info().
			Str("method", info.FullMethod).
			Dur("duration_ms", duration)

		// Handle the error, if any
		if err != nil {
			// Attempt to extract status from gRPC error
			st, ok := status.FromError(err)
			if !ok {
				// If not a gRPC error, log the raw error
				e.Str("error", err.Error()).Msg("gRPC call failed")
			} else {
				// Log the gRPC-specific error details
				e.
					Str("error", st.Message()).
					Str("code", st.Code().String()).
					Msg("gRPC call failed")
			}
		} else {
			e.Msg("gRPC call succeeded")
		}

		return resp, err
	}
}
