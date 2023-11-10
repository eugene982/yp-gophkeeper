package note

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

type NoteWritter interface {
	NoteWrite(ctx context.Context, data storage.NoteData) error
}

type NoteWritterFunc func(ctx context.Context, data storage.NoteData) error

func (f NoteWritterFunc) NoteWrite(ctx context.Context, data storage.NoteData) error {
	return f(ctx, data)
}

var _ NoteWritter = NoteWritterFunc(nil)

type GRPCWriteHandler func(ctx context.Context, in *pb.NoteWriteRequest) (*empty.Empty, error)

// NewGRPCWriteHandler - функция-конструктор ручки записи заметки
func NewGRPCWriteHandler(w NoteWritter, getUserID handler.GetUserIDFunc, enc crypt.Encryptor) GRPCWriteHandler {
	return func(ctx context.Context, in *pb.NoteWriteRequest) (*empty.Empty, error) {
		var err error

		write := storage.NoteData{
			Name: in.Name,
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

		err = w.NoteWrite(ctx, write)
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
