package binary

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"github.com/eugene982/yp-gophkeeper/internal/storage"

	crypt "github.com/eugene982/yp-gophkeeper/internal/crypto"
)

type BinaryWritter interface {
	BinaryWrite(ctx context.Context, data storage.BinaryData) (int64, error)
}

type BinaryWritterFunc func(ctx context.Context, data storage.BinaryData) (int64, error)

func (f BinaryWritterFunc) BinaryWrite(ctx context.Context, data storage.BinaryData) (int64, error) {
	return f(ctx, data)
}

var _ BinaryWritter = BinaryWritterFunc(nil)

type GRPCWriteHandler func(ctx context.Context, in *pb.BinaryWriteRequest) (*pb.BinaryWriteResponse, error)

// NewGRPCWriteHandler - функция-конструктор ручки записи бинарника
func NewGRPCWriteHandler(w BinaryWritter, getUserID handler.GetUserIDFunc, enc crypt.Encryptor) GRPCWriteHandler {
	return func(ctx context.Context, in *pb.BinaryWriteRequest) (*pb.BinaryWriteResponse, error) {
		var err error

		write := storage.BinaryData{
			Name: in.Name,
			Sise: in.Sise,
		}

		write.UserID, err = getUserID(ctx)
		if err != nil {
			return nil, err
		}

		write.Notes, err = enc.Encrypt([]byte(in.Notes))
		if err != nil {
			logger.Errorf("encrypt notes error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		id, err := w.BinaryWrite(ctx, write)
		if err != nil {
			if errors.Is(err, storage.ErrWriteConflict) {
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
			logger.Errorf("write note error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &pb.BinaryWriteResponse{Id: int32(id)}, nil
	}
}
