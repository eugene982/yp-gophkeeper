// Package postgres реализация хранилища на postres
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

var (
	errUnkmownDataType = errors.New("unknown data type")

	writeQuery = map[string]string{ // запросы записи в бд
		// Пользователь
		"users": `INSERT INTO users 
			(user_id, passwd_hash)
		VALUES(:user_id, :passwd_hash);`,

		// Пароль
		"passwords": `INSERT INTO passwords
			(user_id, name, username, password, notes)
		VALUES(:user_id, :name, :username, :password, :notes);`,

		// Карточки
		"cards": `INSERT INTO cards
			(user_id, name, number, pin, notes)
		VALUES(:user_id, :name, :number, :pin, :notes);`,

		// Заметки
		"notes": `INSERT INTO notes
			(user_id, name, notes)
		VALUES(:user_id, :name, :notes);`,

		// бинарники
		"binaries": `INSERT INTO binaries
			(user_id, name, size, notes, bin_id)
		VALUES(:user_id, :name, :size, :notes, :bin_id);`,
	}

	updateQuery = map[string]string{ // запросы на обновление данных
		"users": `UPDATE users 
		SET user_id=:user_id, passwd_hash=:passwd_hash, update_at=now()   
		WHERE user_id=:user_id;`,

		"passwords": `UPDATE passwords 
		SET user_id=:user_id, name=:name, username=:username, password=:password, notes=:notes, update_at=now()   
		WHERE id=:id;`,

		"cards": `UPDATE cards 
		SET user_id=:user_id, name=:name, number=:number, notes=:notes, update_at=now()   
		WHERE id=:id;`,

		"notes": `UPDATE notes 
		SET user_id=:user_id, name=:name, notes=:notes, update_at=now()  
		WHERE id=:id;`,

		"binaries": `UPDATE binaries 
		SET user_id=:user_id, name=:name, size=:size, notes=:notes, update_at=now()   
		WHERE id=:id;`,
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

// Open - Функция открытия БД
func (p *PgxStore) Open(db *sqlx.DB, migratePath string) error {
	if migratePath != "" {
		driver, err := migrate_postgres.WithInstance(db.DB, &migrate_postgres.Config{})
		if err != nil {
			return err
		}

		absParh, err := filepath.Abs(migratePath)
		if err != nil {
			return err
		}

		m, err := migrate.NewWithDatabaseInstance("file://"+absParh,
			"posrgres", driver)
		if err != nil {
			return err
		}
		if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	p.db = db
	return nil
}

// Close закрытие соединения
func (p *PgxStore) Close() error {
	return p.db.Close()
}

// Ping пинг к базе
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

// ReadList читаем баланс пользователя
func (p *PgxStore) ReadList(ctx context.Context, userID string) (res storage.ListData, err error) {
	query := `
	SELECT
		users.user_id AS user_id,
		COUNT(DISTINCT notes.id) AS notes_count,
		COUNT(DISTINCT cards.id) AS cards_count,
		COUNT(DISTINCT binaries.id) AS binaries_count,
		COUNT(DISTINCT passwords.id) AS passwords_count
	FROM users
		LEFT JOIN notes ON users.user_id = notes.user_id
		LEFT JOIN cards ON users.user_id = cards.user_id
		LEFT JOIN binaries ON users.user_id = binaries.user_id
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

// Passwords //

// PasswordList список наименований паролей
func (p *PgxStore) PasswordList(ctx context.Context, userID string) ([]string, error) {
	return p.namesList(ctx, "passwords", userID)
}

// PasswordWrite запись нового пароля
func (p *PgxStore) PasswordWrite(ctx context.Context, data storage.PasswordData) error {
	return p.Write(ctx, data)
}

// PasswordRead чтение пароля из базы
func (p *PgxStore) PasswordRead(ctx context.Context, userID, name string) (res storage.PasswordData, err error) {
	err = p.readFirstByName(ctx, &res, userID, name)
	return
}

// PasswordDelete удаление пароля
func (p *PgxStore) PasswordDelete(ctx context.Context, userID, name string) error {
	return p.deleteByName(ctx, "passwords", userID, name)
}

// PasswordUpdate обновление пароля
func (p *PgxStore) PasswordUpdate(ctx context.Context, data storage.PasswordData) error {
	return p.Update(ctx, data)
}

// Cards //

// CardList список наименований карт
func (p *PgxStore) CardList(ctx context.Context, userID string) ([]string, error) {
	return p.namesList(ctx, "cards", userID)
}

// CardWrite запись новой карты
func (p *PgxStore) CardWrite(ctx context.Context, data storage.CardData) error {
	return p.Write(ctx, data)
}

// CardRead чтение карты
func (p *PgxStore) CardRead(ctx context.Context, userID, name string) (res storage.CardData, err error) {
	err = p.readFirstByName(ctx, &res, userID, name)
	return
}

// CardDelete удаление карты
func (p *PgxStore) CardDelete(ctx context.Context, userID, name string) (err error) {
	err = p.deleteByName(ctx, "cards", userID, name)
	return
}

// CardUpdate обновление сведений
func (p *PgxStore) CardUpdate(ctx context.Context, data storage.CardData) error {
	return p.Update(ctx, data)
}

// Notes //

// NoteList список наименовайний заметок
func (p *PgxStore) NoteList(ctx context.Context, userID string) ([]string, error) {
	return p.namesList(ctx, "notes", userID)
}

// NoteWrite запись новой заметки в базу
func (p *PgxStore) NoteWrite(ctx context.Context, data storage.NoteData) error {
	return p.Write(ctx, data)
}

// NoteRead чтение заметки
func (p *PgxStore) NoteRead(ctx context.Context, userID, name string) (res storage.NoteData, err error) {
	err = p.readFirstByName(ctx, &res, userID, name)
	return
}

// NoteDelete удаление заметки
func (p *PgxStore) NoteDelete(ctx context.Context, userID, name string) (err error) {
	err = p.deleteByName(ctx, "notes", userID, name)
	return
}

// NoteUpdate обновление заметки
func (p *PgxStore) NoteUpdate(ctx context.Context, data storage.NoteData) error {
	return p.Update(ctx, data)
}

// Binary //

// BinaryList список наименований двоичных данных
func (p *PgxStore) BinaryList(ctx context.Context, userID string) ([]string, error) {
	return p.namesList(ctx, "binaries", userID)
}

// BinaryWrite запись нового бинарника, возвращаем идентификатор бинарника.
func (p *PgxStore) BinaryWrite(ctx context.Context, data storage.BinaryData) (int64, error) {
	var err error

	// Запросим новый идентификатор большого объекта
	if data.BinID == 0 {
		query := `SELECT lo_creat(-1)`
		err = p.db.GetContext(ctx, &data.BinID, query)
		if err != nil {
			return 0, err
		}
	}

	err = p.Write(ctx, data)
	if err != nil {
		return 0, errNoContent(err)
	}

	return data.BinID, errNoContent(err)
}

// BinaryRead чтение бинарных данных
func (p *PgxStore) BinaryRead(ctx context.Context, userID, name string) (res storage.BinaryData, err error) {
	err = p.readFirstByName(ctx, &res, userID, name)
	return
}

// BinaryDelete удаление бинарника
func (p *PgxStore) BinaryDelete(ctx context.Context, userID, name string) (err error) {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			logger.Error(fmt.Errorf("psql rollbacck error: %w", err))
		}
	}()

	query := `SELECT lo_unlink(bin_id) 
		FROM binaries WHERE user_id=$1 AND name=$2`

	_, err = tx.ExecContext(ctx, query, userID, name)
	if err != nil {
		return err
	}
	err = p.deleteByName(ctx, "binaries", userID, name)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// BinaryUpdate обновление бинарника
func (p *PgxStore) BinaryUpdate(ctx context.Context, data storage.BinaryData) error {
	var (
		err error
		tx  *sqlx.Tx
	)

	if data.BinID != 0 {
		tx, err = p.db.BeginTxx(ctx, nil)
		if err != nil {
			return err
		}
		defer func() {
			if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
				logger.Error(fmt.Errorf("psql rollbacck error: %w", err))
			}
		}()

		_, err = tx.ExecContext(ctx, "SELECT lo_unlink($1)", data.BinID)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, "SELECT lo_create($1)", data.BinID)
		if err != nil {
			return err
		}
	}

	err = p.Update(ctx, data)
	if err != nil {
		return err
	}
	if tx != nil {
		return tx.Commit()
	}
	return nil
}

// BinaryUpload запись фрагмента бинарника в хранилище бинарника
func (p *PgxStore) BinaryUpload(ctx context.Context, data storage.BinaryChunk) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			logger.Error(fmt.Errorf("psql rollbacck error: %w", err))
		}
	}()

	query := `SELECT lo_put(:bin_id, :offset, :chunk);`

	_, err = tx.NamedExecContext(ctx, query, data)
	if err != nil {
		return errNoContent(err)
	}
	return tx.Commit()
}

// BinaryDownload обновление бинарника
func (p *PgxStore) BinaryDownload(ctx context.Context, data *storage.BinaryChunk) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			logger.Error(fmt.Errorf("psql rollbacck error: %w", err))
		}
	}()

	query := `SELECT lo_get($1, $2, $3);`

	err = tx.GetContext(ctx, &data.Chunk, query, data.BinID, data.Offset, cap(data.Chunk))
	if err != nil {
		return errNoContent(err)
	}
	return tx.Commit()
}

//

func (p *PgxStore) Write(ctx context.Context, data any) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			logger.Error(fmt.Errorf("psql rollbacck error: %w", err))
		}
	}()

	var query string
	switch data.(type) {
	case storage.UserData:
		query = writeQuery["users"]
	case storage.PasswordData:
		query = writeQuery["passwords"]
	case storage.CardData:
		query = writeQuery["cards"]
	case storage.NoteData:
		query = writeQuery["notes"]
	case storage.BinaryData:
		query = writeQuery["binaries"]
	default:
		return errUnkmownDataType
	}

	_, err = tx.NamedExecContext(ctx, query, data)
	if err != nil {
		return errWriteConflict(err)
	}

	return tx.Commit()
}

func (p *PgxStore) Update(ctx context.Context, data any) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			logger.Error(fmt.Errorf("psql rollbacck error: %w", err))
		}
	}()

	var query string
	switch data.(type) {
	case storage.UserData:
		query = updateQuery["users"]
	case storage.PasswordData:
		query = updateQuery["passwords"]
	case storage.CardData:
		query = updateQuery["cards"]
	case storage.NoteData:
		query = updateQuery["notes"]
	case storage.BinaryData:
		query = updateQuery["binaries"]
	default:
		return errUnkmownDataType
	}

	if _, err = tx.NamedExecContext(ctx, query, data); err != nil {
		return errNoContent(errWriteConflict(err))
	}
	return tx.Commit()
}

func (p *PgxStore) deleteByName(ctx context.Context, tabname, userID, name string) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			logger.Error(fmt.Errorf("psql rollbacck error: %w", err))
		}
	}()

	query := `DELETE FROM ` + tabname +
		` WHERE user_id=$1 AND name=$2;`

	res, err := tx.ExecContext(ctx, query, userID, name)
	if err != nil {
		return err
	}
	if n, e := res.RowsAffected(); n == 0 && e == nil {
		return storage.ErrNoContent
	}
	return tx.Commit()
}

func (p *PgxStore) namesList(ctx context.Context, tabname, userID string) ([]string, error) {
	query := `SELECT name FROM ` + tabname +
		` WHERE user_id = $1`

	res := make([]string, 0)
	err := p.db.SelectContext(ctx, &res, query, userID)

	if err != nil {
		err = errNoContent(err)
		return nil, err
	}
	return res, nil
}

func (p *PgxStore) readFirstByName(ctx context.Context, res any, userID, name string) error {
	var tabname string

	switch res.(type) {
	case *storage.CardData:
		tabname = "cards"
	case *storage.NoteData:
		tabname = "notes"
	case *storage.PasswordData:
		tabname = "passwords"
	case *storage.BinaryData:
		tabname = "binaries"
	default:
		return errUnkmownDataType
	}

	query := `SELECT * FROM ` + tabname +
		` WHERE user_id = $1 AND name = $2 LIMIT 1`

	err := p.db.GetContext(ctx, res, query, userID, name)
	return errNoContent(err)
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
