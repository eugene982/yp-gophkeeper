package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ErrWriteConflict = errors.New("write conflict")
	ErrNoContent     = errors.New("no content")
)

var database Storage

func Open(dns string) (Storage, error) {

	if dns == "" {
		return nil, errors.New("database dsn is empty")
	}

	db, err := sqlx.Open("pgx", dns)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(3 * time.Minute)

	if database == nil {
		return nil, errors.New("database not initialize")
	}

	err = database.Open(db)
	if err != nil {
		return nil, err
	}

	return database, nil
}

// регистрируем подключенный драйвер
func RegDriver(db Storage) {
	database = db
}

// Интерфейс для хранилища данных
type Storage interface {
	Open(*sqlx.DB) error
	Close() error
	Ping(context.Context) error
	WriteUser(context.Context, UserData) error
	ReadUser(context.Context, string) (UserData, error)
	ReadList(context.Context, string) (ListData, error)

	//Password
	PasswordList(context.Context, string) ([]string, error)
	PasswordWrite(context.Context, PasswordData) error
	PasswordRead(context.Context, string, string) (PasswordData, error)
	PasswordUpdate(context.Context, PasswordData) error
	PasswordDelete(context.Context, string, string) error
}
