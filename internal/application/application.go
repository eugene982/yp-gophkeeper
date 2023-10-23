package application

import (
	"github.com/eugene982/yp-gophkeeper/internal/config"
	crypt "github.com/eugene982/yp-gophkeeper/internal/crypto"
	aescrypt "github.com/eugene982/yp-gophkeeper/internal/crypto/aes"
	"github.com/eugene982/yp-gophkeeper/internal/grpc"

	"github.com/eugene982/yp-gophkeeper/internal/storage"
	_ "github.com/eugene982/yp-gophkeeper/internal/storage/postgres"
)

type Application struct {
	grpcServer *grpc.GRPCServer
	storage    storage.Storage
	crypt      crypt.EncryptDecryptor
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

	app.crypt, err = aescrypt.New(grpc.CRYPTO_KEY)
	if err != nil {
		return nil, err
	}

	app.grpcServer, err = grpc.NewServer(app.storage, app.crypt, conf.ServerAddres)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

// Start запуск прослушивания
func (app *Application) Start() error {
	return app.grpcServer.Start()
}
