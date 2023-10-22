package password

import (
	"context"
	"errors"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PasswordUpdater interface {
	PasswordUpdate(ctx context.Context, data storage.PasswordData) error
}

type PasswordUpdateFunc func(ctx context.Context, data storage.PasswordData) error

func (f PasswordUpdateFunc) PasswordUpdate(ctx context.Context, data storage.PasswordData) error {
	return f(ctx, data)
}

var _ PasswordUpdater = PasswordUpdateFunc(nil)

type GRPCUpdateHandler func(context.Context, *pb.PasswordUpdateRequest) (*empty.Empty, error)

func NewGRPCUpdateHandler(u PasswordUpdater, getUserID handler.GetUserIDFunc) GRPCUpdateHandler {
	return func(ctx context.Context, in *pb.PasswordUpdateRequest) (*empty.Empty, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		data := fromPasswordWriteRequest(userID, in.Id, in.Write)
		err = u.PasswordUpdate(ctx, data)
		if err != nil {
			if errors.Is(err, storage.ErrWriteConflict) {
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
		}

		return &empty.Empty{}, nil
	}
}
