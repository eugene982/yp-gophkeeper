package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHash - хеширование пароля
func PasswordHash(password, salt string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash - сверка пароля с хешем
func CheckPasswordHash(hash, password, salt string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+salt))
	return err == nil
}
