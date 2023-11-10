package password

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

type PasswordWritter interface {
	PasswordWrite(ctx context.Context, data storage.PasswordData) error
}

type PasswordWritterFunc func(ctx context.Context, data storage.PasswordData) error

func (f PasswordWritterFunc) PasswordWrite(ctx context.Context, data storage.PasswordData) error {
	return f(ctx, data)
}

var _ PasswordWritter = PasswordWritterFunc(nil)

type GRPCWriteHandler func(ctx context.Context, in *pb.PasswordWriteRequest) (*empty.Empty, error)

// NewGRPCWriteHandler - функция-конструктор ручки записи пароля
func NewGRPCWriteHandler(w PasswordWritter, getUserID handler.GetUserIDFunc, enc crypt.Encryptor) GRPCWriteHandler {
	return func(ctx context.Context, in *pb.PasswordWriteRequest) (*empty.Empty, error) {
		var err error

		write := storage.PasswordData{
			Name: in.Name,
		}

		write.UserID, err = getUserID(ctx)
		if err != nil {
			return nil, err
		}

		write.Username, err = enc.Encrypt([]byte(in.Username))
		if err != nil {
			logger.Errorf("encrypt username error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		write.Password, err = enc.Encrypt([]byte(in.Password))
		if err != nil {
			logger.Errorf("encrypt password error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		write.Notes, err = enc.Encrypt([]byte(in.Notes))
		if err != nil {
			logger.Errorf("encrypt notes error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		err = w.PasswordWrite(ctx, write)
		if err != nil {
			if errors.Is(err, storage.ErrWriteConflict) {
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
			logger.Errorf("write password error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &empty.Empty{}, nil
	}
}
