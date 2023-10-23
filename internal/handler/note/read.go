package note

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

type NoteReader interface {
	NoteRead(ctx context.Context, userID, name string) (storage.NoteData, error)
}

type NoteReaderFunc func(ctx context.Context, userID, name string) (storage.NoteData, error)

func (f NoteReaderFunc) NoteRead(ctx context.Context, userID, name string) (storage.NoteData, error) {
	return f(ctx, userID, name)
}

var _ NoteReader = NoteReaderFunc(nil)

type GRPCReadHandler func(context.Context, *pb.NoteReadRequest) (*pb.NoteReadResponse, error)

func NewGRPCReadHandler(r NoteReader, getUserID handler.GetUserIDFunc, dec crypt.Decryptor) GRPCReadHandler {
	return func(ctx context.Context, in *pb.NoteReadRequest) (*pb.NoteReadResponse, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		data, err := r.NoteRead(ctx, userID, in.Name)
		if err != nil {
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			logger.Errorf("read note error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		notes, err := dec.Decrypt(data.Notes)
		if err != nil {
			logger.Errorf("decrypt notes error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		resp := pb.NoteReadResponse{
			Id:    data.ID,
			Name:  data.Name,
			Notes: string(notes),
		}

		return &resp, nil
	}
}
