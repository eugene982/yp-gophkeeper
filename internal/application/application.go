// Package application реализация приложения
package application

import (
	"github.com/eugene982/yp-gophkeeper/internal/config"
	crypt "github.com/eugene982/yp-gophkeeper/internal/crypto"
	aescrypt "github.com/eugene982/yp-gophkeeper/internal/crypto/aes"
	grpc_v1 "github.com/eugene982/yp-gophkeeper/internal/grpc/v1"

	"github.com/eugene982/yp-gophkeeper/internal/storage"
	_ "github.com/eugene982/yp-gophkeeper/internal/storage/postgres"
)

// Application структура основного приложения
type Application struct {
	grpcServer *grpc_v1.GRPCServer
	storage    storage.Storage
	crypt      crypt.EncryptDecryptor
}

// New конструктор
func New(conf config.Config) (*Application, error) {
	var (
		app Application
		err error
	)

	app.storage, err = storage.Open(conf.DSN, conf.MigratePath)
	if err != nil {
		return nil, err
	}

	app.crypt, err = aescrypt.New(grpc_v1.CryptoKey)
	if err != nil {
		return nil, err
	}

	app.grpcServer, err = grpc_v1.NewServer(app.storage, app.crypt, conf.ServerAddres)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

// Start запуск прослушивания
func (app *Application) Start() error {
	return app.grpcServer.Start()
}

// Stop остановка приложения
func (app *Application) Stop() error {
	app.grpcServer.Stop()
	return app.storage.Close()
}
