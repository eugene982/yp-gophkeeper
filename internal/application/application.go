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

// // Register регистрация пользователя
// func (app *Application) Register(ctx context.Context, login, password string) (token string, err error) {
// 	defer func() { // залогируем ошибки при появлении
// 		if err != nil && err != handler.ErrUnauthenticated {
// 			logger.Errorf("register error: %w", err)
// 		}
// 	}()

// 	hash, err := utils.PasswordHash(password, PASSWORD_SALT)
// 	if err != nil {
// 		return
// 	}

// 	err = app.storage.WriteUser(ctx, storage.UserData{
// 		UserID:       login,
// 		PasswordHash: hash,
// 	})
// 	if err != nil {
// 		if errors.Is(err, storage.ErrWriteConflict) {
// 			err = handler.ErrAlreadyExists
// 		}
// 		return
// 	}

// 	token, err = jwt.MakeToken(login, TOKEN_SECRET_KEY, TOKEN_EXP)
// 	return
// }

// func (app *Application) Login(ctx context.Context, login, password string) (token string, err error) {
// 	defer func() { // залогируем ошибки при появлении
// 		if err != nil && err != handler.ErrUnauthenticated {
// 			logger.Errorf("login error: %w", err)
// 		}
// 	}()

// 	data, err := app.storage.ReadUser(ctx, login)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrNoContent) {
// 			err = handler.ErrUnauthenticated
// 		}
// 		return
// 	}

// 	if !utils.CheckPasswordHash(data.PasswordHash, password, PASSWORD_SALT) {
// 		err = handler.ErrUnauthenticated
// 		return
// 	}

// 	token, err = jwt.MakeToken(login, TOKEN_SECRET_KEY, TOKEN_EXP)
// 	return
// }

// func (app *Application) List(ctx context.Context, userID string) (storage.ListData, error) {
// 	res, err := app.storage.ReadList(ctx, userID)
// 	if err != nil {
// 		logger.Errorf("login error: %w", err,
// 			"user_id", userID)
// 	}
// 	return res, err
// }

// func (app *Application) NamesList(ctx context.Context, tab storage.TableName, userID string) ([]string, error) {
// 	// res, err := app.storage.NamesList(ctx, tab, userID)
// 	// if err != nil {
// 	// 	logger.Errorf("login error: %w", err,
// 	// 		"table", tab,
// 	// 		"user_id", userID)
// 	// }
// 	// return res, err
// }

//unc (app *Application) Write(ctx context.Context, data any) error {
// 	if err := app.storage.Write(ctx, data); err != nil {
// 		logger.Errorf("write error: %w", err)
// 		return nil
// 	}
// 	return nil
// }

// func (app *Application) ReadByName(ctx context.Context, tab storage.TableName, userID, name string) (any, error) {
// 	res, err := app.storage.ReadByName(ctx, tab, userID, name)
// 	if err != nil {
// 		logger.Errorf("read error: %w", err,
// 			"table", tab,
// 			"user_id", userID,
// 			"name", name)
// 		return nil, err
// 	}
// 	return res, nil
// }

// func (app *Application) DeleteByName(ctx context.Context, tab storage.TableName, userID, name string) error {
// 	err := app.storage.DeleteByName(ctx, tab, userID, name)
// 	if err != nil {
// 		logger.Errorf("delete error: %w", err,
// 			"table", tab,
// 			"user_id", userID,
// 			"name", name)
// 		return err
// 	}
// 	return nil
// }
