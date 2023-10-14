package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims - структура утверждений, которая включает стандартные утверждения
// и одно пользовательское UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// BuildJWTString cоздаёт токен и возвращает его в виде строки.
func BuildJWTString(userID string, secret_key string, exp time.Duration) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждением - Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
		// собственное утверждение
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secret_key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetJWTUserID возвращает идентификатор пользователя
func GetJWTUserID(tokenString, secret_key string) (string, error) {
	// создаём экземпляр утверждения
	claims := &Claims{}
	// парсим из строки токена tokenString в структуру
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok { // Проверка на совпадения метода подписи
			return nil, fmt.Errorf("unexpected signed method: %v", t.Header["alg"])
		}

		return []byte(secret_key), nil
	})

	if err != nil {
		return "", nil
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	// возвращаем ID полезователя
	return claims.UserID, nil
}
