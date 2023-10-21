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

type CardWritter interface {
	CardWrite(ctx context.Context, data storage.CardData) error
}

type CardWritterFunc func(ctx context.Context, data storage.CardData) error

func (f CardWritterFunc) CardWrite(ctx context.Context, data storage.CardData) error {
	return f(ctx, data)
}

var _ CardWritter = CardWritterFunc(nil)

type GRPCWriteHandler func(ctx context.Context, in *pb.CardWriteRequest) (*empty.Empty, error)

func NewGRPCWriteHandler(w CardWritter, getUserID handler.GetUserIDFunc) GRPCWriteHandler {
	return func(ctx context.Context, in *pb.CardWriteRequest) (*empty.Empty, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		err = w.CardWrite(ctx, fromCardWriteRequest(userID, in))
		if err != nil {
			if errors.Is(err, storage.ErrWriteConflict) {
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
			logger.Errorf("card note error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &empty.Empty{}, nil
	}
}
