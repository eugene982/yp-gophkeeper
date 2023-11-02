package card

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

type CardUpdater interface {
	CardUpdate(ctx context.Context, data storage.CardData) error
}

type CardUpdaterFunc func(ctx context.Context, data storage.CardData) error

func (f CardUpdaterFunc) CardUpdate(ctx context.Context, data storage.CardData) error {
	return f(ctx, data)
}

var _ CardUpdater = CardUpdaterFunc(nil)

type GRPCUpdateHandler func(context.Context, *pb.CardUpdateRequest) (*empty.Empty, error)

// NewGRPCUpdateHandler - функция конструктор ручки для обновления карты
func NewGRPCUpdateHandler(u CardUpdater, getUserID handler.GetUserIDFunc, enc crypt.Encryptor) GRPCUpdateHandler {
	return func(ctx context.Context, in *pb.CardUpdateRequest) (*empty.Empty, error) {
		var err error

		upd := storage.CardData{
			ID:   in.Id,
			Name: in.Write.Name,
		}

		upd.UserID, err = getUserID(ctx)
		if err != nil {
			return nil, err
		}

		upd.Number, err = enc.Encrypt([]byte(in.Write.Number))
		if err != nil {
			logger.Errorf("encrypt number error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		upd.Pin, err = enc.Encrypt([]byte(in.Write.Pin))
		if err != nil {
			logger.Errorf("encrypt pin error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		upd.Notes, err = enc.Encrypt([]byte(in.Write.Notes))
		if err != nil {
			logger.Errorf("encrypt notes error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		err = u.CardUpdate(ctx, upd)
		if err != nil {
			if errors.Is(err, storage.ErrWriteConflict) {
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			logger.Errorf("card update error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &empty.Empty{}, nil
	}
}
