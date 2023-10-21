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

type PasswordWritter interface {
	PasswordWrite(ctx context.Context, data storage.PasswordData) error
}

type PasswordWritterFunc func(ctx context.Context, data storage.PasswordData) error

func (f PasswordWritterFunc) PasswordWrite(ctx context.Context, data storage.PasswordData) error {
	return f(ctx, data)
}

var _ PasswordWritter = PasswordWritterFunc(nil)

type GRPCWriteHandler func(ctx context.Context, in *pb.PasswordWriteRequest) (*empty.Empty, error)

func NewGRPCWriteHandler(w PasswordWritter, getUserID handler.GetUserIDFunc) GRPCWriteHandler {
	return func(ctx context.Context, in *pb.PasswordWriteRequest) (*empty.Empty, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		err = w.PasswordWrite(ctx, fromPasswordWriteRequest(userID, in))
		if err != nil {
			if errors.Is(err, storage.ErrWriteConflict) {
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
			logger.Errorf("write password error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &empty.Empty{}, nil
	}
}
