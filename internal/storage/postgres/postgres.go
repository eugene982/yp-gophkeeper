// Хранение в базе данных postres
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

type PgxStore struct {
	db *sqlx.DB
}

func init() {
	storage.RegDriver(new(PgxStore))
}

// Утверждение типа, ошибка компиляции
var _ storage.Storage = (*PgxStore)(nil)

// Функция открытия БД
func (p *PgxStore) Open(db *sqlx.DB) error {

	if err := createTablesIfNonExists(db); err != nil {
		return err
	}
	p.db = db
	return nil
}

// Закрытие соединения
func (p *PgxStore) Close() error {
	return p.db.Close()
}

// Пинг к базе
func (p *PgxStore) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

// WriteUser Установка уникального соответствия
func (p *PgxStore) WriteUser(ctx context.Context, data storage.UserData) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO users (user_id, passwd_hash) 
		VALUES(:user_id, :passwd_hash);`

	if _, err = tx.NamedExecContext(ctx, query, data); err != nil {
		return errWriteConflict(err)
	}
	return tx.Commit()
}

// ReadUser Чтение данных пользователя
func (p *PgxStore) ReadUser(ctx context.Context, userID string) (res storage.UserData, err error) {
	query := `
		SELECT * FROM users
		WHERE user_id = $1 LIMIT 1`

	if err = p.db.GetContext(ctx, &res, query, userID); err != nil {
		err = errNoContent(err)
	}
	return
}

// При первом запуске база может быть пустая
func createTablesIfNonExists(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			user_id VARCHAR (100) PRIMARY KEY,
			passwd_hash TEXT NOT NULL
		);
		`
	_, err := db.Exec(query)
	return err
}

func errWriteConflict(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		return storage.ErrWriteConflict
	}
	return err
}

func errNoContent(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrNoContent
	}
	return err
}
