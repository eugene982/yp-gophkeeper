package application

import (
	"context"
	"fmt"

	"github.com/eugene982/yp-gophkeeper/internal/config"
	"github.com/eugene982/yp-gophkeeper/internal/grpc"
	"github.com/eugene982/yp-gophkeeper/internal/utils"
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

// Start запуск прослушивания
func (app *Application) Start() error {
	return app.grpcServer.Start()
}

// Ping проверка соединения
func (app *Application) Ping(ctx context.Context) error {
	// заглушка
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// Register регистрация пользователя
func (app *Application) Register(ctx context.Context, login, password string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	token, err := utils.BuildJWTString(login, grpc.SECRET_KEY, grpc.TOKEN_EXP)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (app *Application) Login(ctx context.Context, login, passwd string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return "", fmt.Errorf("not implements")
	}
}
