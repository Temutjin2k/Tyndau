package main

import (
	"context"

	"github.com/Temutjin2k/Tyndau/api-gateway/config"
	"github.com/Temutjin2k/Tyndau/api-gateway/internal/app"
	"github.com/Temutjin2k/Tyndau/api-gateway/pkg/zerologger"
)

func main() {
	ctx := context.Background()

	logger := zerologger.NewZeroLogger()

	cfg, err := config.New()
	if err != nil {
		logger.Error().Err(err)
		return
	}

	app, err := app.New(ctx, cfg, logger)
	if err != nil {
		return
	}

	err = app.Run()
	if err != nil {
		return
	}
}
