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
)

type NoteDeleter interface {
	NoteDelete(ctx context.Context, userID, name string) error
}

type NoteDeleteFunc func(ctx context.Context, userID, name string) error

func (f NoteDeleteFunc) NoteDelete(ctx context.Context, userID, name string) error {
	return f(ctx, userID, name)
}

var _ NoteDeleter = NoteDeleteFunc(nil)

type GRPCDeleteHandler func(ctx context.Context, in *pb.NoteDelRequest) (*empty.Empty, error)

// NewGRPCDeleteHandler - функция-конструктор ручки удаления заметки
func NewGRPCDeleteHandler(d NoteDeleter, getUserID handler.GetUserIDFunc) GRPCDeleteHandler {
	return func(ctx context.Context, in *pb.NoteDelRequest) (*empty.Empty, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		err = d.NoteDelete(ctx, userID, in.Name)
		if err != nil {
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			logger.Errorf("delete note error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &empty.Empty{}, nil
	}
}
