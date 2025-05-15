package app

import (
	"context"
	"fmt"
	"github.com/Temutjin2k/Tyndau/music-service/internal/song/cache/redis"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Temutjin2k/Tyndau/music-service/config"
	grpcserver "github.com/Temutjin2k/Tyndau/music-service/internal/adapter/grpc/server"
	postgresSong "github.com/Temutjin2k/Tyndau/music-service/internal/song/adapter/postgres"
	"github.com/Temutjin2k/Tyndau/music-service/internal/song/usecase"
	postgresDB "github.com/Temutjin2k/Tyndau/music-service/pkg/postgres"
	"github.com/rs/zerolog"
)

const serviceName = "music-service"

type App struct {
	grpcServer *grpcserver.API
	postgresDB *postgresDB.PostgreDB
	logger     *zerolog.Logger
}

func New(ctx context.Context, cfg *config.Config, logger *zerolog.Logger) (*App, error) {
	postgresDB, err := postgresDB.New(ctx, cfg.Postgres)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to Postgres")
		return nil, fmt.Errorf("postgres: %w", err)
	}

	ttlSec, _ := strconv.Atoi(os.Getenv("CACHE_TTL_SECONDS"))
	if ttlSec == 0 {
		ttlSec = 300 // default fallback
	}
	cacheSvc := redis.New(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), 0, time.Duration(ttlSec)*time.Second)

	repo := postgresSong.NewSongRepo(postgresDB.Pool)
	songService := usecase.NewSongService(repo, cacheSvc)

	grpcSrv := grpcserver.New(cfg.Server.GRPCServer, logger, songService)

	app := &App{
		grpcServer: grpcSrv,
		postgresDB: postgresDB,
		logger:     logger,
	}

	return app, nil
}

func (a *App) Close(ctx context.Context) {
	err := a.grpcServer.Stop(ctx)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to stop grpc server")
	}
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	ctx := context.Background()

	a.grpcServer.Run(ctx, errCh)

	a.logger.Info().Str("name", serviceName).Msg("service started")

	// Wait for shutdown
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return errRun
	case s := <-shutdownCh:
		a.logger.Info().Str("signal", s.String()).Msg("shutting down app")
		a.Close(ctx)
	}

	return nil
}
