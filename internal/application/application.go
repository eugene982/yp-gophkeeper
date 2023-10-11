package application

import (
	"context"

	"github.com/eugene982/yp-gophkeeper/internal/config"
)

type Application struct{}

func New(config.Config) (*Application, error) {
	var app Application
	return &app, nil
}

func (app *Application) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return nil
}
