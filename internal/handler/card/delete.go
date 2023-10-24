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

type CardDeleter interface {
	CardDelete(ctx context.Context, userID, name string) error
}

type CardDeleteFunc func(ctx context.Context, userID, name string) error

func (f CardDeleteFunc) CardDelete(ctx context.Context, userID, name string) error {
	return f(ctx, userID, name)
}

var _ CardDeleter = CardDeleteFunc(nil)

type GRPCDeleteHandler func(ctx context.Context, in *pb.CardDelRequest) (*empty.Empty, error)

func NewGRPCDeleteHandler(d CardDeleter, getUserID handler.GetUserIDFunc) GRPCDeleteHandler {
	return func(ctx context.Context, in *pb.CardDelRequest) (*empty.Empty, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		err = d.CardDelete(ctx, userID, in.Name)
		if err != nil {
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			logger.Errorf("card delete error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &empty.Empty{}, nil
	}
}
