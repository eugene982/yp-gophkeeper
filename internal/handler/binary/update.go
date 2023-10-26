package binary

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

type BinaryUpdater interface {
	BinaryUpdate(ctx context.Context, data storage.BinaryData) error
}

type BinaryUpdaterFunc func(ctx context.Context, data storage.BinaryData) error

func (f BinaryUpdaterFunc) BinaryUpdate(ctx context.Context, data storage.BinaryData) error {
	return f(ctx, data)
}

var _ BinaryUpdater = BinaryUpdaterFunc(nil)

type GRPCUpdateHandler func(context.Context, *pb.BinaryUpdateRequest) (*empty.Empty, error)

func NewGRPCUpdateHandler(u BinaryUpdater, getUserID handler.GetUserIDFunc, enc crypt.Encryptor) GRPCUpdateHandler {
	return func(ctx context.Context, in *pb.BinaryUpdateRequest) (*empty.Empty, error) {
		var err error

		upd := storage.BinaryData{
			ID:   in.Id,
			Name: in.Write.Name,
		}

		upd.UserID, err = getUserID(ctx)
		if err != nil {
			return nil, err
		}

		upd.Bin, err = enc.Encrypt(in.Write.Bin)
		if err != nil {
			logger.Errorf("encrypt bin error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		upd.Notes, err = enc.Encrypt([]byte(in.Write.Notes))
		if err != nil {
			logger.Errorf("encrypt notes error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		err = u.BinaryUpdate(ctx, upd)
		if err != nil {
			if errors.Is(err, storage.ErrWriteConflict) {
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			logger.Errorf("binary update error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &empty.Empty{}, nil
	}
}
