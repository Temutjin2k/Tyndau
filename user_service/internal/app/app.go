package app

import (
	"context"

	"github.com/Temutjin2k/Tyndau/user_svc/config"
)

type Application struct {
}

func New(ctx context.Context, config *config.Config) (*Application, error) {
	return &Application{}, nil
}

func (app *Application) Run() error {
	return nil
}
