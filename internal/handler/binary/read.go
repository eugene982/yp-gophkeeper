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

type BinaryReader interface {
	BinaryRead(ctx context.Context, userID, name string) (storage.BinaryData, error)
}

type BinaryReaderFunc func(ctx context.Context, userID, name string) (storage.BinaryData, error)

func (f BinaryReaderFunc) BinaryRead(ctx context.Context, userID, name string) (storage.BinaryData, error) {
	return f(ctx, userID, name)
}

var _ BinaryReader = BinaryReaderFunc(nil)

type GRPCReadHandler func(context.Context, *pb.BinaryReadRequest) (*pb.BinaryReadResponse, error)

func NewGRPCReadHandler(r BinaryReader, getUserID handler.GetUserIDFunc, dec crypt.Decryptor) GRPCReadHandler {
	return func(ctx context.Context, in *pb.BinaryReadRequest) (*pb.BinaryReadResponse, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		data, err := r.BinaryRead(ctx, userID, in.Name)
		if err != nil {
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			logger.Errorf("read note error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		bin, err := dec.Decrypt(data.Bin)
		if err != nil {
			logger.Errorf("decrypt bin error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		notes, err := dec.Decrypt(data.Notes)
		if err != nil {
			logger.Errorf("decrypt notes error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		resp := pb.BinaryReadResponse{
			Id:    data.ID,
			Name:  data.Name,
			Bin:   bin,
			Notes: string(notes),
		}

		return &resp, nil
	}
}
