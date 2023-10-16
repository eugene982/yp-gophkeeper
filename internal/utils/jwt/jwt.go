package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims - структура утверждений, которая включает стандартные утверждения
// и одно пользовательское UserID
type claims struct {
	jwt.RegisteredClaims
	UserID string
}

// MakeToken cоздаёт токен и возвращает его в виде строки.
func MakeToken(userID string, secret_key string, exp time.Duration) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждением - Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
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

// GetUserID возвращает идентификатор пользователя
func GetUserID(token, secret_key string) (string, error) {
	// создаём экземпляр утверждения
	claims := &claims{}
	// парсим из строки токена tokenString в структуру
	jwtoken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok { // Проверка на совпадения метода подписи
			return nil, fmt.Errorf("unexpected signed method: %v", t.Header["alg"])
		}

		return []byte(secret_key), nil
	})

	if err != nil {
		return "", nil
	}

	if !jwtoken.Valid {
		return "", fmt.Errorf("invalid token")
	}
	// возвращаем ID полезователя
	return claims.UserID, nil
}
