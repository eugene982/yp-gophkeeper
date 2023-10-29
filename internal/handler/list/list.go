package list

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

// List интерфейс возвращающий количество данных пользователя
type ListReader interface {
	ReadList(context.Context, string) (storage.ListData, error)
}

type ListReaderFunc func(ctx context.Context, userID string) (storage.ListData, error)

func (f ListReaderFunc) ReadList(ctx context.Context, userID string) (storage.ListData, error) {
	return f(ctx, userID)
}

var _ ListReader = ListReaderFunc(nil)

type GRPCHandler func(context.Context, *empty.Empty) (*pb.ListResponse, error)

func NewRPCListHandler(list ListReader, getUserID handler.GetUserIDFunc) GRPCHandler {
	return func(ctx context.Context, _ *empty.Empty) (*pb.ListResponse, error) {

		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		data, err := list.ReadList(ctx, userID)
		if err != nil {
			logger.Errorf("read user list error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &pb.ListResponse{
			CardsCount:     int32(data.CardsCount),
			NotesCount:     int32(data.NotesCount),
			BinariesCount:  int32(data.BinariesCount),
			PasswordsCount: int32(data.PasswordsCount),
		}, nil
	}
}
