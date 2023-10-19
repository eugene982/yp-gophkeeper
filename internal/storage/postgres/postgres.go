// Хранение в базе данных postres
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

const (
	// текст запроса создания таблиц
	createQuery = `
		CREATE TABLE IF NOT EXISTS users (
			user_id     VARCHAR(64) PRIMARY KEY,
			passwd_hash TEXT        NOT NULL
		);

		CREATE TABLE IF NOT EXISTS passwords (
			user_id     VARCHAR(64)  NOT NULL,
			name	    VARCHAR(128) NOT NULL,			
			username    VARCHAR(128) NOT NULL,
			password    VARCHAR(128) NOT NULL,
			notes       TEXT         NOT NULL
		);
		CREATE INDEX IF NOT EXISTS user_id_idx 
		ON passwords (user_id);
		CREATE UNIQUE INDEX IF NOT EXISTS user_id_name_idx 
		ON passwords (user_id, name);	


		CREATE TABLE IF NOT EXISTS notes (
			user_id VARCHAR(64)  NOT NULL,
			name	VARCHAR(128) NOT NULL,
			notes   TEXT         NOT NULL
		);
		CREATE INDEX IF NOT EXISTS user_id_idx 
		ON notes (user_id);
		CREATE UNIQUE INDEX IF NOT EXISTS user_id_name_idx 
		ON notes (user_id, name);	
		
		CREATE TABLE IF NOT EXISTS cards (
			user_id VARCHAR(64)  NOT NULL,
			name	VARCHAR(128) NOT NULL,
			number	VARCHAR(20)  NOT NULL,
			notes   TEXT         NOT NULL
		);
		CREATE INDEX IF NOT EXISTS user_id_idx 
		ON cards (user_id);
		CREATE UNIQUE INDEX IF NOT EXISTS user_id_name_idx 
		ON cards (user_id, name);
		`
)

var (
	errUnkmownDataType = errors.New("unknown data type")

	writeQuery = map[string]string{ // запросы записи в бд
		// Пользователь
		"users": `INSERT INTO users (user_id, passwd_hash)
		VALUES(:user_id, :passwd_hash);`,

		// Пароль
		"passwords": `INSERT INTO passwords (user_id, name, username, password, notes)
		VALUES(:user_id, :name, :username, :password, :notes);`,
	}

	updateQuery = map[string]string{ // запросы на обновление данных
		"users":     "",
		"passwords": "",
	}
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
	return p.Write(ctx, data)
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

// читаем баланс пользователя
func (p *PgxStore) ReadList(ctx context.Context, userID string) (res storage.ListData, err error) {
	query := `
	SELECT
		users.user_id AS user_id,
		COUNT(notes.name) AS notes_count,
		COUNT(cards.name) AS cards_count,
		COUNT(passwords.name) AS passwords_count
	FROM users
		LEFT JOIN notes ON users.user_id = notes.user_id
		LEFT JOIN cards ON cards.user_id = notes.user_id
		LEFT JOIN passwords ON users.user_id = passwords.user_id
	WHERE 
		users.user_id = $1
	GROUP BY 
		users.user_id;`

	if err = p.db.GetContext(ctx, &res, query, userID); err != nil {
		err = errNoContent(err)
	}
	return
}

func (p *PgxStore) namesList(ctx context.Context, tabname storage.TableName, userID string) ([]string, error) {
	query := fmt.Sprintf(`
		SELECT name FROM %s
		WHERE user_id = $1`, tabname)

	res := make([]string, 0)
	err := p.db.GetContext(ctx, &res, query, userID)

	if err != nil {
		err = errNoContent(err)
		return nil, err
	}
	return res, nil
}

func (p *PgxStore) PasswordList(ctx context.Context, userID string) ([]string, error) {
	return p.namesList(ctx, "passwords", userID)
}

func (p *PgxStore) CardList(ctx context.Context, userID string) ([]string, error) {
	return p.namesList(ctx, "cards", userID)
}

func (p *PgxStore) NoteList(ctx context.Context, userID string) ([]string, error) {
	return p.namesList(ctx, "notes", userID)
}

func (p *PgxStore) Write(ctx context.Context, data any) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var query string
	switch data.(type) {
	case storage.UserData:
		query = writeQuery["users"]
	case storage.PasswordData:
		query = writeQuery["passwords"]
	default:
		return errUnkmownDataType
	}

	if _, err = tx.NamedExecContext(ctx, query, data); err != nil {
		return errWriteConflict(err)
	}
	return tx.Commit()
}

func (p *PgxStore) Update(ctx context.Context, data any) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var query string
	switch data.(type) {
	case storage.UserData:
		query = updateQuery["users"]
	case storage.PasswordData:
		query = updateQuery["passwords"]
	default:
		return errUnkmownDataType
	}

	if _, err = tx.NamedExecContext(ctx, query, data); err != nil {
		return errWriteConflict(err)
	}
	return tx.Commit()
}

func (p *PgxStore) ReadByName(ctx context.Context, tabname storage.TableName, userID, name string) (any, error) {
	var res any

	switch tabname {
	case "passwords":
		res = storage.PasswordData{}
	default:
		return nil, errUnkmownDataType
	}

	query := `SELECT * FROM ` + string(tabname) +
		` WHERE user_id = $1 AND name = $2 LIMIT 1`

	err := p.db.GetContext(ctx, &res, query, userID, name)
	return res, errNoContent(err)
}

// При первом запуске база может быть пустая
func createTablesIfNonExists(db *sqlx.DB) error {
	_, err := db.Exec(createQuery)
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
