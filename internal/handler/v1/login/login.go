package login

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

type UserReader interface {
	ReadUser(context.Context, string) (storage.UserData, error)
}

type UserReaderFunc func(context.Context, string) (storage.UserData, error)

func (f UserReaderFunc) ReadUser(ctx context.Context, userId string) (storage.UserData, error) {
	return f(ctx, userId)
}

var _ UserReader = UserReaderFunc(nil)

type GRPCHandler func(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error)

type HashCheckFunc func(string, string) bool
type TokenGenFunc func(string) (string, error)

// NewRPCLoginHandler
func NewRPCLoginHandler(r UserReader, checkFn HashCheckFunc, tokenFn TokenGenFunc) GRPCHandler {
	return func(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {

		data, err := r.ReadUser(ctx, in.Login)
		if err != nil {
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.Unauthenticated, err.Error())
			}
			logger.Errorf("read user error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		if !checkFn(in.Password, data.PasswordHash) {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		var resp pb.LoginResponse
		resp.Token, err = tokenFn(in.Login)
		if err != nil {
			logger.Errorf("make token error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &resp, nil
	}
}
