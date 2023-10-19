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

// Специальный тип чтоб не перепутать параметры
type TableName string

// Интерфейс для хранилища данных
type Storage interface {
	Open(*sqlx.DB) error
	Close() error
	Ping(context.Context) error
	WriteUser(context.Context, UserData) error
	ReadUser(context.Context, string) (UserData, error)
	ReadList(context.Context, string) (ListData, error)

	// Password
	PasswordList(ctx context.Context, userID string) ([]string, error)

	// Card
	CardList(ctx context.Context, userID string) ([]string, error)

	// Notes
	NoteList(ctx context.Context, userID string) ([]string, error)

	//Write(ctx context.Context, data any) error
	//Update(ctx context.Context, data any) error

	//NamesList(ctx context.Context, tab TableName, userID string) ([]string, error)
	//ReadByName(ctx context.Context, tab TableName, userID, name string) (any, error)
	//DeleteByName(ctx context.Context, tab TableName, userID, name string) error
}
