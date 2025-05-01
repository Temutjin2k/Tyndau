package main

import (
	"context"
	"os"

	"github.com/Temutjin2k/Tyndau/user_service/config"
	"github.com/Temutjin2k/Tyndau/user_service/internal/app"
	"github.com/Temutjin2k/Tyndau/user_service/pkg/zerologger"
)

func main() {
	ctx := context.Background()

	logger := zerologger.NewZeroLogger()

	// Parse config
	cfg, err := config.New()
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse config")
		os.Exit(1)
	}

	application, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Error().Err(err).Msg("failed to setup application")
		os.Exit(1)
	}

	err = application.Run()
	if err != nil {
		logger.Error().Err(err).Msg("failed to run application")
		os.Exit(1)
	}
}
