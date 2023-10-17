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
	UserID   string `db:"user_id"`
	Name     string `db:"name"`
	Username string `db:"username"`
	Password string `db:"password"`
	Notes    string `db:"notes"`
}
