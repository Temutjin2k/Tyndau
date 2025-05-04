package server

import (
	"context"
	"fmt"
	"net"

	"github.com/Temutjin2k/Tyndau/music-service/config"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type API struct {
	server *grpc.Server
	cfg    config.GRPCServer
	addr   string
	logger *zerolog.Logger
}

func New(cfg config.GRPCServer, logger *zerolog.Logger) *API {
	addr := fmt.Sprintf("127.0.0.1:%d", cfg.Port)

	return &API{
		cfg:    cfg,
		addr:   addr,
		logger: logger,
	}
}

func (a *API) Run(ctx context.Context, errCh chan<- error) {
	go func() {
		a.logger.Info().
			Str("address", a.addr).
			Msg("gRPC server starting")

		if err := a.run(ctx); err != nil {
			a.logger.Error().
				Err(err).
				Msg("gRPC server failed to start")
			errCh <- fmt.Errorf("can't start grpc server: %w", err)
		}
	}()
}

func (a *API) Stop(ctx context.Context) error {
	if a.server == nil {
		a.logger.Debug().Msg("gRPC server already stopped (nil server instance)")
		return nil
	}

	stopped := make(chan struct{})
	go func() {
		a.logger.Info().Msg("gRPC server shutting down gracefully")
		a.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		a.logger.Warn().Msg("gRPC server force stopped due to context cancellation")
		a.server.Stop()
	case <-stopped:
		a.logger.Info().Msg("gRPC server stopped gracefully")
	}

	return nil
}

func (a *API) run(ctx context.Context) error {
	a.server = grpc.NewServer(a.setOptions(ctx)...)

	// Register services
	reflection.Register(a.server)

	a.logger.Debug().Msg("gRPC services registered")

	listener, err := net.Listen("tcp", a.addr)
	if err != nil {
		a.logger.Error().
			Err(err).
			Str("address", a.addr).
			Msg("failed to create listener")
		return fmt.Errorf("failed to create listener: %w", err)
	}

	a.logger.Info().
		Str("address", a.addr).
		Msg("gRPC server started listening")

	if err := a.server.Serve(listener); err != nil {
		a.logger.Error().
			Err(err).
			Msg("gRPC server failed to serve")
		return fmt.Errorf("failed to serve grpc: %w", err)
	}

	return nil
}
