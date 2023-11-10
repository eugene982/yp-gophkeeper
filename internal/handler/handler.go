// Package handler хранит общие функции
package handler

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrRPCInvalidToken = status.Errorf(codes.Unauthenticated, "invalid token")
)

type GetUserIDFunc func(context.Context) (string, error)

// GetUserIDFromMD функция получает идентификатор пользователя из токена
func GetUserIDFromMD(ctx context.Context, secretKey string) (string, error) {
	var token string

	if md, ok := metadata.FromIncomingContext(ctx); !ok {
		return "", ErrRPCInvalidToken
	} else if vals := md.Get("token"); len(vals) > 0 && vals[0] != "" {
		token = vals[0]
	} else {
		return "", ErrRPCInvalidToken
	}

	userID, err := GetUserID(token, secretKey)
	if err != nil {
		return "", ErrRPCInvalidToken
	}
	return userID, nil
}

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

// Claims - структура утверждений, которая включает стандартные утверждения
// и одно пользовательское UserID
type claims struct {
	jwt.RegisteredClaims
	UserID string
}

// MakeToken cоздаёт токен и возвращает его в виде строки.
func MakeToken(userID string, secretKey string, exp time.Duration) (string, error) {
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
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetUserID возвращает идентификатор пользователя
func GetUserID(token, secretKey string) (string, error) {
	// создаём экземпляр утверждения
	claims := &claims{}
	// парсим из строки токена tokenString в структуру
	jwtoken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok { // Проверка на совпадения метода подписи
			return nil, fmt.Errorf("unexpected signed method: %v", t.Header["alg"])
		}

		return []byte(secretKey), nil
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
