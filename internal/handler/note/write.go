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

type NoteWritter interface {
	NoteWrite(ctx context.Context, data storage.NoteData) error
}

type NoteWritterFunc func(ctx context.Context, data storage.NoteData) error

func (f NoteWritterFunc) NoteWrite(ctx context.Context, data storage.NoteData) error {
	return f(ctx, data)
}

var _ NoteWritter = NoteWritterFunc(nil)

type GRPCWriteHandler func(ctx context.Context, in *pb.NoteWriteRequest) (*empty.Empty, error)

func NewGRPCWriteHandler(w NoteWritter, getUserID handler.GetUserIDFunc) GRPCWriteHandler {
	return func(ctx context.Context, in *pb.NoteWriteRequest) (*empty.Empty, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		err = w.NoteWrite(ctx, fromNoteWriteRequest(userID, in))
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
