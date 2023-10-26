package password

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"github.com/eugene982/yp-gophkeeper/internal/storage"

	crypt "github.com/eugene982/yp-gophkeeper/internal/crypto"
)

type PasswordUpdater interface {
	PasswordUpdate(ctx context.Context, data storage.PasswordData) error
}

type PasswordUpdaterFunc func(ctx context.Context, data storage.PasswordData) error

func (f PasswordUpdaterFunc) PasswordUpdate(ctx context.Context, data storage.PasswordData) error {
	return f(ctx, data)
}

var _ PasswordUpdater = PasswordUpdaterFunc(nil)

type GRPCUpdateHandler func(context.Context, *pb.PasswordUpdateRequest) (*empty.Empty, error)

func NewGRPCUpdateHandler(u PasswordUpdater, getUserID handler.GetUserIDFunc, enc crypt.Encryptor) GRPCUpdateHandler {
	return func(ctx context.Context, in *pb.PasswordUpdateRequest) (*empty.Empty, error) {
		var err error

		upd := storage.PasswordData{
			ID:   in.Id,
			Name: in.Write.Name,
		}

		upd.UserID, err = getUserID(ctx)
		if err != nil {
			return nil, err
		}

		upd.Username, err = enc.Encrypt([]byte(in.Write.Username))
		if err != nil {
			logger.Errorf("encrypt username error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		upd.Password, err = enc.Encrypt([]byte(in.Write.Password))
		if err != nil {
			logger.Errorf("encrypt password error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		upd.Notes, err = enc.Encrypt([]byte(in.Write.Notes))
		if err != nil {
			logger.Errorf("encrypt notes error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		err = u.PasswordUpdate(ctx, upd)
		if err != nil {
			if errors.Is(err, storage.ErrWriteConflict) {
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			logger.Errorf("update password error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &empty.Empty{}, nil
	}
}
