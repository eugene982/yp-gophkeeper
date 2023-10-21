package card

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

type CardListGetter interface {
	CardList(ctx context.Context, userID string) ([]string, error)
}

type CardListGetterFunc func(ctx context.Context, userID string) ([]string, error)

func (f CardListGetterFunc) CardList(ctx context.Context, userID string) ([]string, error) {
	return f(ctx, userID)
}

var _ CardListGetter = CardListGetterFunc(nil)

type GRPCListHandler func(ctx context.Context, in *empty.Empty) (*pb.CardListResponse, error)

func NewGRPCListHandler(g CardListGetter, getUserID handler.GetUserIDFunc) GRPCListHandler {
	return func(ctx context.Context, in *empty.Empty) (*pb.CardListResponse, error) {

		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		var resp pb.CardListResponse
		resp.Names, err = g.CardList(ctx, userID)
		if err != nil {
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}

			logger.Errorf("read card list error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &resp, nil
	}
}
