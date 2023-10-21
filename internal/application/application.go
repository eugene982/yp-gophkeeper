package application

import (
	"github.com/eugene982/yp-gophkeeper/internal/config"
	"github.com/eugene982/yp-gophkeeper/internal/grpc"

	"github.com/eugene982/yp-gophkeeper/internal/storage"
	_ "github.com/eugene982/yp-gophkeeper/internal/storage/postgres"
)

type Application struct {
	grpcServer *grpc.GRPCServer
	storage    storage.Storage
}

// New конструктор
func New(conf config.Config) (*Application, error) {
	var (
		app Application
		err error
	)

	app.storage, err = storage.Open(conf.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	app.grpcServer, err = grpc.NewServer(app.storage, conf.ServerAddres)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

// Start запуск прослушивания
func (app *Application) Start() error {
	return app.grpcServer.Start()
}
