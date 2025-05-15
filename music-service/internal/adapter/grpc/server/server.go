package server

import (
	"context"
	"fmt"
	"net"

	"github.com/Temutjin2k/Tyndau/music-service/config"
	"github.com/Temutjin2k/Tyndau/music-service/internal/adapter/grpc/server/frontend"
	pb "github.com/Temutjin2k/Tyndau/music-service/internal/api/song/v1"
	"github.com/Temutjin2k/Tyndau/music-service/internal/song/usecase"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type API struct {
	server *grpc.Server
	cfg    config.GRPCServer
	addr   string
	logger *zerolog.Logger

	songUC usecase.SongService
}

func New(cfg config.GRPCServer, log *zerolog.Logger, uc usecase.SongService) *API {
	return &API{
		cfg:    cfg,
		addr:   fmt.Sprintf("127.0.0.1:%d", cfg.Port),
		logger: log,
		songUC: uc,
	}
}

func (a *API) Run(ctx context.Context, errCh chan<- error) {
	go func() {
		a.logger.Info().Str("addr", a.addr).Msg("gRPC starting")
		if err := a.run(ctx); err != nil {
			a.logger.Error().Err(err).Msg("gRPC failed")
			errCh <- err
		}
	}()
}

func (a *API) Stop(ctx context.Context) error {
	if a.server == nil {
		return nil
	}
	stopped := make(chan struct{})
	go func() {
		a.logger.Info().Msg("gRPC shutting down")
		a.server.GracefulStop()
		close(stopped)
	}()
	select {
	case <-ctx.Done():
		a.logger.Warn().Msg("force stop gRPC")
		a.server.Stop()
	case <-stopped:
	}
	return nil
}

/* ------------ private ------------ */

func (a *API) run(ctx context.Context) error {
	a.server = grpc.NewServer(a.setOptions(ctx)...)

	// регистрируем Song-хэндлер
	pb.RegisterSongServiceServer(a.server, frontend.NewSongHandler(a.songUC, a.logger))
	reflection.Register(a.server)

	lis, err := net.Listen("tcp", a.addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", a.addr, err)
	}
	return a.server.Serve(lis)
}
