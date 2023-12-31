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

type NoteUpdater interface {
	NoteUpdate(ctx context.Context, data storage.NoteData) error
}

type NoteUpdaterFunc func(ctx context.Context, data storage.NoteData) error

func (f NoteUpdaterFunc) NoteUpdate(ctx context.Context, data storage.NoteData) error {
	return f(ctx, data)
}

var _ NoteUpdater = NoteUpdaterFunc(nil)

type GRPCUpdateHandler func(context.Context, *pb.NoteUpdateRequest) (*empty.Empty, error)

// NewGRPCUpdateHandler - функция конструктор ручки для обновления заметки
func NewGRPCUpdateHandler(u NoteUpdater, getUserID handler.GetUserIDFunc, enc crypt.Encryptor) GRPCUpdateHandler {
	return func(ctx context.Context, in *pb.NoteUpdateRequest) (*empty.Empty, error) {
		var err error

		upd := storage.NoteData{
			ID:   in.Id,
			Name: in.Write.Name,
		}

		upd.UserID, err = getUserID(ctx)
		if err != nil {
			return nil, err
		}

		upd.Notes, err = enc.Encrypt([]byte(in.Write.Notes))
		if err != nil {
			logger.Errorf("encrypt notes error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		err = u.NoteUpdate(ctx, upd)
		if err != nil {
			if errors.Is(err, storage.ErrWriteConflict) {
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			logger.Errorf("note update error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &empty.Empty{}, nil
	}
}
