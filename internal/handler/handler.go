package handler

import (
	"context"
	"errors"

	"github.com/eugene982/yp-gophkeeper/internal/utils/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrAlreadyExists   = errors.New("already exists")
	ErrUnauthenticated = errors.New("unauthenticated")

	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
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
		return "", errMissingMetadata
	} else if vals := md.Get("token"); len(vals) > 0 {
		token = vals[0]
	} else {
		return "", errInvalidToken
	}

	userID, err := jwt.GetUserID(token, secret_key)
	if err != nil {
		return "", errInvalidToken
	}
	return userID, nil
}
