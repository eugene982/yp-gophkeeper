package handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/eugene982/yp-gophkeeper/internal/utils/jwt"
)

var (
	ErrAlreadyExists   = errors.New("already exists")
	ErrUnauthenticated = errors.New("unauthenticated")

	ErrRPCInvalidToken = status.Errorf(codes.Unauthenticated, "invalid token")
)

type UserIDGetter interface {
	GetUserID(context.Context) (string, error)
}

type UserIDGetterFunc func(context.Context) (string, error)

func (f UserIDGetterFunc) GetUserID(ctx context.Context) (string, error) {
	return f(ctx)
}

func NewMDUserIDGetter(secret_key string) UserIDGetter {
	return UserIDGetterFunc(func(ctx context.Context) (string, error) {
		return GetUserIDFromMD(ctx, secret_key)
	})
}

func GetUserIDFromMD(ctx context.Context, secret_key string) (string, error) {
	var token string

	if md, ok := metadata.FromIncomingContext(ctx); !ok {
		return "", ErrRPCInvalidToken
	} else if vals := md.Get("token"); len(vals) > 0 && vals[0] != "" {
		token = vals[0]
	} else {
		return "", ErrRPCInvalidToken
	}

	userID, err := jwt.GetUserID(token, secret_key)
	if err != nil {
		return "", ErrRPCInvalidToken
	}
	return userID, nil
}
