package card

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"github.com/eugene982/yp-gophkeeper/internal/storage"

	crypt "github.com/eugene982/yp-gophkeeper/internal/crypto"
)

type CardReader interface {
	CardRead(ctx context.Context, userID, name string) (storage.CardData, error)
}

type CardReaderFunc func(ctx context.Context, userID, name string) (storage.CardData, error)

func (f CardReaderFunc) CardRead(ctx context.Context, userID, name string) (storage.CardData, error) {
	return f(ctx, userID, name)
}

var _ CardReader = CardReaderFunc(nil)

type GRPCReadHandler func(context.Context, *pb.CardReadRequest) (*pb.CardReadResponse, error)

func NewGRPCReadHandler(r CardReader, getUserID handler.GetUserIDFunc, dec crypt.Decryptor) GRPCReadHandler {
	return func(ctx context.Context, in *pb.CardReadRequest) (*pb.CardReadResponse, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		data, err := r.CardRead(ctx, userID, in.Name)
		if err != nil {
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			logger.Errorf("read card error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		number, err := dec.Decrypt(data.Number)
		if err != nil {
			logger.Errorf("decrypt numder error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		pin, err := dec.Decrypt(data.Pin)
		if err != nil {
			logger.Errorf("decrypt pin error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		notes, err := dec.Decrypt(data.Notes)
		if err != nil {
			logger.Errorf("decrypt notes error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		resp := pb.CardReadResponse{
			Id:     data.ID,
			Name:   data.Name,
			Number: string(number),
			Pin:    string(pin),
			Notes:  string(notes),
		}

		return &resp, nil
	}
}
