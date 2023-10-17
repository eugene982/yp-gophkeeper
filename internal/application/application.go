package application

import (
	"context"
	"errors"
	"time"

	"github.com/eugene982/yp-gophkeeper/internal/config"
	"github.com/eugene982/yp-gophkeeper/internal/grpc"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"github.com/eugene982/yp-gophkeeper/internal/utils"
	"github.com/eugene982/yp-gophkeeper/internal/utils/jwt"

	"github.com/eugene982/yp-gophkeeper/internal/storage"
	_ "github.com/eugene982/yp-gophkeeper/internal/storage/postgres"
)

var (
	TOKEN_SECRET_KEY = "sekret=key"
	TOKEN_EXP        = time.Hour
	PASSWORD_SALT    = "password=salt"
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

	app.grpcServer, err = grpc.NewServer(&app, conf.ServerAddres, TOKEN_SECRET_KEY)
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
	return app.storage.Ping(ctx)
}

// Register регистрация пользователя
func (app *Application) Register(ctx context.Context, login, password string) (token string, err error) {
	defer func() { // залогируем ошибки при появлении
		if err != nil && err != handler.ErrUnauthenticated {
			logger.Errorf("register error: %w", err)
		}
	}()

	hash, err := utils.PasswordHash(password, PASSWORD_SALT)
	if err != nil {
		return
	}

	err = app.storage.WriteUser(ctx, storage.UserData{
		UserID:       login,
		PasswordHash: hash,
	})
	if err != nil {
		if errors.Is(err, storage.ErrWriteConflict) {
			err = handler.ErrAlreadyExists
		}
		return
	}

	token, err = jwt.MakeToken(login, TOKEN_SECRET_KEY, TOKEN_EXP)
	return
}

func (app *Application) Login(ctx context.Context, login, password string) (token string, err error) {
	defer func() { // залогируем ошибки при появлении
		if err != nil && err != handler.ErrUnauthenticated {
			logger.Errorf("login error: %w", err)
		}
	}()

	data, err := app.storage.ReadUser(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrNoContent) {
			err = handler.ErrUnauthenticated
		}
		return
	}

	if !utils.CheckPasswordHash(data.PasswordHash, password, PASSWORD_SALT) {
		err = handler.ErrUnauthenticated
		return
	}

	token, err = jwt.MakeToken(login, TOKEN_SECRET_KEY, TOKEN_EXP)
	return
}

func (app *Application) List(ctx context.Context, userID string) (storage.ListData, error) {
	res, err := app.storage.ReadList(ctx, userID)
	if err != nil {
		logger.Errorf("login error: %w", err)
	}
	return res, err
}

func (app *Application) PasswordList(context.Context, string) ([]string, error) {
	panic("not implement")
}

func (app *Application) PasswordWrite(context.Context, storage.PasswordData) error {
	panic("not implement")
}

func (app *Application) PasswordRead(context.Context, string, string) (storage.PasswordData, error) {
	panic("not implement")
}

func (app *Application) PasswordUpdate(context.Context, storage.PasswordData) error {
	panic("not implement")
}

func (app *Application) PasswordDelete(context.Context, string, string) error {
	panic("not implement")
}
