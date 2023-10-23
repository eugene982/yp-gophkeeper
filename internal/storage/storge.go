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

	// Password
	PasswordList(ctx context.Context, userID string) ([]string, error)
	PasswordRead(ctx context.Context, userID, name string) (PasswordData, error)
	PasswordWrite(ctx context.Context, data PasswordData) error
	PasswordDelete(ctx context.Context, userID, name string) error
	PasswordUpdate(ctx context.Context, data PasswordData) error

	// Card
	CardList(ctx context.Context, userID string) ([]string, error)
	CardRead(ctx context.Context, userID, name string) (CardData, error)
	CardWrite(ctx context.Context, data CardData) error
	CardDelete(ctx context.Context, userID, name string) error
	CardUpdate(ctx context.Context, data CardData) error

	// Notes
	NoteList(ctx context.Context, userID string) ([]string, error)
	NoteRead(ctx context.Context, userID, name string) (NoteData, error)
	NoteWrite(ctx context.Context, data NoteData) error
	NoteDelete(ctx context.Context, userID, name string) error
	NoteUpdate(ctx context.Context, data NoteData) error

	// Binary
	BinaryList(ctx context.Context, userID string) ([]string, error)
	BinaryRead(ctx context.Context, userID, name string) (BinaryData, error)
	BinaryWrite(ctx context.Context, data BinaryData) error
	BinaryDelete(ctx context.Context, userID, name string) error
	BinaryUpdate(ctx context.Context, data BinaryData) error
}
