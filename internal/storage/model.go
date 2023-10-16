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
	PasswordsCount int `db:"passwords_count"`
	NotesCount     int `db:"notes_count"`
	CardsCount     int `db:"cards_count"`
}
