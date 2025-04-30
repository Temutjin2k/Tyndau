package main

import (
	"context"
	"os"

	"github.com/Temutjin2k/Tyndau/api-gateway/config"
	"github.com/Temutjin2k/Tyndau/api-gateway/internal/app"
	"github.com/Temutjin2k/Tyndau/api-gateway/pkg/zerologger"
)

func main() {
	// Context
	ctx := context.Background()

	// Logger
	logger := zerologger.NewZeroLogger()

	cfg, err := config.New()
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse config")
		os.Exit(1)
	}

	app, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Error().Err(err).Msg("failed to setup application")
		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		logger.Error().Err(err).Msg("failed to run application")
		os.Exit(1)
	}
}
