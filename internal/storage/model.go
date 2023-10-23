package storage

// ListData структура хранилища содержащая информацию о количестве
// хранимой информации по пользователю
type UserData struct {
	UserID       string `db:"user_id"`
	PasswordHash string `db:"passwd_hash"`
}

// ListData структура хранилища содержащая информацию о количестве
// хранимой информации по пользователю
type ListData struct {
	UserID         string `db:"user_id"`
	NotesCount     int    `db:"notes_count"`
	CardsCount     int    `db:"cards_count"`
	PasswordsCount int    `db:"passwords_count"`
}

// PasswordData хранимая информация о паролях
type PasswordData struct {
	ID       int32  `db:"id"`
	UserID   string `db:"user_id"`
	Name     string `db:"name"`
	Username []byte `db:"username"`
	Password []byte `db:"password"`
	Notes    []byte `db:"notes"`
}

type CardData struct {
	ID     int32  `db:"id"`
	UserID string `db:"user_id"`
	Name   string `db:"name"`
	Number []byte `db:"number"`
	Pin    []byte `db:"pin"`
	Notes  []byte `db:"notes"`
}

type NoteData struct {
	ID     int32  `db:"id"`
	UserID string `db:"user_id"`
	Name   string `db:"name"`
	Notes  []byte `db:"notes"`
}

type BinaryData struct {
	ID     int32  `db:"id"`
	UserID string `db:"user_id"`
	Name   string `db:"name"`
	Bin    []byte `db:"bin"`
	Notes  []byte `db:"notes"`
}
