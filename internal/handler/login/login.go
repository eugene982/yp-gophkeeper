package login

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
)

// Login интерфейс отвечающий за авторизацию пользователей
type Login interface {
	Login(ctx context.Context, login, passwd string) (string, error)
}

type LoginFunc func(ctx context.Context, login, passwd string) (string, error)

func (f LoginFunc) Login(ctx context.Context, login, passwd string) (string, error) {
	return f(ctx, login, passwd)
}

var _ Login = LoginFunc(nil)

type GRPCHandler func(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error)

// NewRPCLoginHandler
func NewRPCLoginHandler(login Login) GRPCHandler {
	return func(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
		var (
			resp pb.LoginResponse
			err  error
		)

		resp.Token, err = login.Login(ctx, in.Login, in.Password)
		if err == nil {
			return &resp, nil
		} else if errors.Is(err, handler.ErrUnauthenticated) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
}
