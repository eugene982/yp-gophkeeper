package register

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
)

// Register интерфейс отвечающий за регистрацию пользователей
type Register interface {
	Register(ctx context.Context, login, passwd string) (string, error)
}

type RegisterFunc func(ctx context.Context, login, passwd string) (string, error)

func (f RegisterFunc) Register(ctx context.Context, login, passwd string) (string, error) {
	return f(ctx, login, passwd)
}

var _ Register = RegisterFunc(nil)

type GRPCHandler func(context.Context, *pb.RegisterRequest) (*pb.RegisterResponse, error)

// NewRPCRegisterHandler - ручка регистрации нового пользователя
func NewRPCRegisterHandler(register Register) GRPCHandler {
	return func(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
		var (
			resp pb.RegisterResponse
			err  error
		)

		resp.Token, err = register.Register(ctx, in.Login, in.Password)
		if err == nil {
			return &resp, nil
		} else if errors.Is(err, handler.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
}
