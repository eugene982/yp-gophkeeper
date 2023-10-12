package application

import (
	"context"

	"github.com/eugene982/yp-gophkeeper/internal/config"
	"github.com/eugene982/yp-gophkeeper/internal/grpc"
)

type Application struct {
	grpcServer *grpc.GRPCServer
}

// New конструктор
func New(conf config.Config) (*Application, error) {
	var (
		app Application
		err error
	)

	app.grpcServer, err = grpc.NewServer(&app, conf.ServerAddres)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

// запуск прослушивания
func (app *Application) Start() error {
	return app.grpcServer.Start()
}

func (app *Application) Ping(ctx context.Context) error {
	// заглушка
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (app *Application) Register(ctx context.Context, login, password string) (string, error) {
	// заглушка
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return "secret-token", nil
	}
}
