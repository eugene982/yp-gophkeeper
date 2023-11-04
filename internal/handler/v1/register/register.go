// Package register ручка регистрации пользователя
package register

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

// UserWriter интерфейс отвечающий за регистрацию пользователей
type UserWriter interface {
	WriteUser(context.Context, storage.UserData) error
}

type UserWriterFunc func(context.Context, storage.UserData) error

func (f UserWriterFunc) WriteUser(ctx context.Context, data storage.UserData) error {
	return f(ctx, data)
}

var _ UserWriter = UserWriterFunc(nil)

type GRPCHandler func(context.Context, *pb.RegisterRequest) (*pb.RegisterResponse, error)

type PasswordHashFunc func(string) (string, error)
type TokenGenFunc func(string) (string, error)

// NewRPCRegisterHandler - ручка регистрации нового пользователя
func NewRPCRegisterHandler(w UserWriter, hashFn PasswordHashFunc, tokenFn TokenGenFunc) GRPCHandler {
	return func(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {

		hash, err := hashFn(in.Password)
		if err != nil {
			logger.Errorf("password hash error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		var resp pb.RegisterResponse

		resp.Token, err = tokenFn(in.Login)
		if err != nil {
			logger.Errorf("make token error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		err = w.WriteUser(ctx, storage.UserData{
			UserID:       in.Login,
			PasswordHash: hash,
		})

		if err == nil {
			return &resp, nil
		} else if errors.Is(err, storage.ErrWriteConflict) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
}
