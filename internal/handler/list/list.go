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
type ListGetter interface {
	List(ctx context.Context, userID string) (storage.ListData, error)
}

type ListGetterFunc func(ctx context.Context, userID string) (storage.ListData, error)

func (f ListGetterFunc) List(ctx context.Context, userID string) (storage.ListData, error) {
	return f(ctx, userID)
}

var _ ListGetter = ListGetterFunc(nil)

type GRPCHandler func(context.Context, *empty.Empty) (*pb.ListResponse, error)

func NewRPCListHandler(list ListGetter, ug handler.UserIDGetter) GRPCHandler {
	return func(ctx context.Context, _ *empty.Empty) (*pb.ListResponse, error) {

		userID, err := ug.GetUserID(ctx)
		if err != nil {
			return nil, err
		}

		data, err := list.List(ctx, userID)
		if err != nil {
			logger.Errorf("error rpc list handler: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		resp := pb.ListResponse{
			PasswordsCount: int32(data.PasswordsCount),
			NotesCount:     int32(data.NotesCount),
			CardsCount:     int32(data.CardsCount),
		}
		return &resp, nil
	}
}
