package note

import (
	"context"
	"errors"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NoteListGetter interface {
	NoteList(ctx context.Context, userID string) ([]string, error)
}

type NoteListGetterFunc func(ctx context.Context, userID string) ([]string, error)

func (f NoteListGetterFunc) NoteList(ctx context.Context, userID string) ([]string, error) {
	return f(ctx, userID)
}

var _ NoteListGetter = NoteListGetterFunc(nil)

type GRPCListHandler func(ctx context.Context, in *empty.Empty) (*pb.NoteListResponse, error)

// NewGRPCListHandler - функция-конструктор ручки получения списка наименований заметок
func NewGRPCListHandler(g NoteListGetter, getUserID handler.GetUserIDFunc) GRPCListHandler {
	return func(ctx context.Context, in *empty.Empty) (*pb.NoteListResponse, error) {

		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		var resp pb.NoteListResponse
		resp.Names, err = g.NoteList(ctx, userID)
		if err != nil {
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}

			logger.Errorf("read notes list error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &resp, nil
	}
}
