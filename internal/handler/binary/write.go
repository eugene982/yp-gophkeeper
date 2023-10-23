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

type BinaryWritter interface {
	BinaryWrite(ctx context.Context, data storage.BinaryData) error
}

type BinaryWritterFunc func(ctx context.Context, data storage.BinaryData) error

func (f BinaryWritterFunc) BinaryWrite(ctx context.Context, data storage.BinaryData) error {
	return f(ctx, data)
}

var _ BinaryWritter = BinaryWritterFunc(nil)

type GRPCWriteHandler func(ctx context.Context, in *pb.BinaryWriteRequest) (*empty.Empty, error)

func NewGRPCWriteHandler(w BinaryWritter, getUserID handler.GetUserIDFunc, enc crypt.Encryptor) GRPCWriteHandler {
	return func(ctx context.Context, in *pb.BinaryWriteRequest) (*empty.Empty, error) {
		var err error

		write := storage.BinaryData{
			Name: in.Name,
		}

		write.UserID, err = getUserID(ctx)
		if err != nil {
			return nil, err
		}

		write.Bin, err = enc.Encrypt(in.Bin)
		if err != nil {
			logger.Errorf("encrypt bin error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		write.Notes, err = enc.Encrypt([]byte(in.Notes))
		if err != nil {
			logger.Errorf("encrypt notes error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		err = w.BinaryWrite(ctx, write)
		if err != nil {
			if errors.Is(err, storage.ErrWriteConflict) {
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
			logger.Errorf("write note error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &empty.Empty{}, nil
	}
}
