package password

import (
	"context"
	"errors"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PasswordReader interface {
	PasswordRead(ctx context.Context, userID, name string) (storage.PasswordData, error)
}

type PasswordReaderFunc func(ctx context.Context, userID, name string) (storage.PasswordData, error)

func (f PasswordReaderFunc) PasswordRead(ctx context.Context, userID, name string) (storage.PasswordData, error) {
	return f(ctx, userID, name)
}

type GRPCReadHandler func(context.Context, *pb.PasswordReadRequest) (*pb.PasswordReadResponse, error)

func NewGRPCReadHandler(r PasswordReader, getUserID handler.GetUserIDFunc) GRPCReadHandler {
	return func(ctx context.Context, in *pb.PasswordReadRequest) (*pb.PasswordReadResponse, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		data, err := r.PasswordRead(ctx, userID, in.Name)
		if err != nil {
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
		}

		resp := toPasswordReadResponse(data)
		return &resp, nil
	}
}
