package password

import (
	"context"
	"errors"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PasswordDeleter interface {
	PasswordDelete(ctx context.Context, userID, name string) error
}

type PasswordDeleteFunc func(ctx context.Context, userID, name string) error

func (f PasswordDeleteFunc) PasswordDelete(ctx context.Context, userID, name string) error {
	return f(ctx, userID, name)
}

var _ PasswordDeleter = PasswordDeleteFunc(nil)

type GRPCDeleteHandler func(ctx context.Context, in *pb.PasswordDelRequest) (*empty.Empty, error)

// NewGRPCDeleteHandler - функция-конструктор ручки удаления пароля
func NewGRPCDeleteHandler(d PasswordDeleter, getUserID handler.GetUserIDFunc) GRPCDeleteHandler {
	return func(ctx context.Context, in *pb.PasswordDelRequest) (*empty.Empty, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		err = d.PasswordDelete(ctx, userID, in.Name)
		if err != nil {
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			logger.Errorf("delete password error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &empty.Empty{}, nil
	}
}
