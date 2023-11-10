package storage

import "time"

// UserData структура пользователя
type UserData struct {
	UserID       string    `db:"user_id"`
	PasswordHash string    `db:"passwd_hash"`
	CtreatAt     time.Time `db:"create_at"`
	UpdateAt     time.Time `db:"update_at"`
}

// ListData структура хранилища содержащая информацию о количестве
// хранимой информации по пользователю
type ListData struct {
	UserID         string    `db:"user_id"`
	NotesCount     int       `db:"notes_count"`
	CardsCount     int       `db:"cards_count"`
	BinariesCount  int       `db:"binaries_count"`
	PasswordsCount int       `db:"passwords_count"`
	CtreatAt       time.Time `db:"create_at"`
	UpdateAt       time.Time `db:"update_at"`
}

// PasswordData хранимая информация о паролях
type PasswordData struct {
	ID       int64     `db:"id"`
	UserID   string    `db:"user_id"`
	Name     string    `db:"name"`
	Username []byte    `db:"username"`
	Password []byte    `db:"password"`
	Notes    []byte    `db:"notes"`
	CtreatAt time.Time `db:"create_at"`
	UpdateAt time.Time `db:"update_at"`
}

// CardData информация о различных картах
type CardData struct {
	ID       int64     `db:"id"`
	UserID   string    `db:"user_id"`
	Name     string    `db:"name"`
	Number   []byte    `db:"number"`
	Pin      []byte    `db:"pin"`
	Notes    []byte    `db:"notes"`
	CtreatAt time.Time `db:"create_at"`
	UpdateAt time.Time `db:"update_at"`
}

// NoteData различные заметки, текст
type NoteData struct {
	ID       int64     `db:"id"`
	UserID   string    `db:"user_id"`
	Name     string    `db:"name"`
	Notes    []byte    `db:"notes"`
	CtreatAt time.Time `db:"create_at"`
	UpdateAt time.Time `db:"update_at"`
}

// BinaryData двоичные данные, файлы
type BinaryData struct {
	ID       int64     `db:"id"`
	UserID   string    `db:"user_id"`
	Name     string    `db:"name"`
	Size     int64     `db:"size"`
	Notes    []byte    `db:"notes"`
	BinID    int64     `db:"bin_id"`
	CtreatAt time.Time `db:"create_at"`
	UpdateAt time.Time `db:"update_at"`
}

type BinaryChunk struct {
	BinID  int64  `db:"bin_id"`
	Offset int64  `db:"offset"`
	Chunk  []byte `db:"chunk"`
}
