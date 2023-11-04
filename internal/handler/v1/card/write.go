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

// CardWritter - интерфейс для записи карты
type CardWritter interface {
	CardWrite(ctx context.Context, data storage.CardData) error
}

type CardWritterFunc func(ctx context.Context, data storage.CardData) error

func (f CardWritterFunc) CardWrite(ctx context.Context, data storage.CardData) error {
	return f(ctx, data)
}

var _ CardWritter = CardWritterFunc(nil)

type GRPCWriteHandler func(ctx context.Context, in *pb.CardWriteRequest) (*empty.Empty, error)

// NewGRPCWriteHandler - функция-конструктор ручки записи карты
func NewGRPCWriteHandler(w CardWritter, getUserID handler.GetUserIDFunc, enc crypt.Encryptor) GRPCWriteHandler {
	return func(ctx context.Context, in *pb.CardWriteRequest) (*empty.Empty, error) {
		var err error

		write := storage.CardData{
			Name: in.Name,
		}

		write.UserID, err = getUserID(ctx)
		if err != nil {
			return nil, err
		}

		write.Number, err = enc.Encrypt([]byte(in.Number))
		if err != nil {
			logger.Errorf("encrypt number error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		write.Pin, err = enc.Encrypt([]byte(in.Pin))
		if err != nil {
			logger.Errorf("encrypt pin error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		write.Notes, err = enc.Encrypt([]byte(in.Notes))
		if err != nil {
			logger.Errorf("encrypt notes error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		err = w.CardWrite(ctx, write)
		if err != nil {
			if errors.Is(err, storage.ErrWriteConflict) {
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
			logger.Errorf("write card error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &empty.Empty{}, nil
	}
}
