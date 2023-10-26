package password

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

type PasswordReader interface {
	PasswordRead(ctx context.Context, userID, name string) (storage.PasswordData, error)
}

type PasswordReaderFunc func(ctx context.Context, userID, name string) (storage.PasswordData, error)

func (f PasswordReaderFunc) PasswordRead(ctx context.Context, userID, name string) (storage.PasswordData, error) {
	return f(ctx, userID, name)
}

var _ PasswordReader = PasswordReaderFunc(nil)

type GRPCReadHandler func(context.Context, *pb.PasswordReadRequest) (*pb.PasswordReadResponse, error)

func NewGRPCReadHandler(r PasswordReader, getUserID handler.GetUserIDFunc, dec crypt.Decryptor) GRPCReadHandler {
	return func(ctx context.Context, in *pb.PasswordReadRequest) (*pb.PasswordReadResponse, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		data, err := r.PasswordRead(ctx, userID, in.Name)
		if err != nil {
			if errors.Is(err, storage.ErrNoContent) {
				return nil, status.Error(codes.NotFound, err.Error())
			}
			logger.Errorf("read password error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		username, err := dec.Decrypt(data.Username)
		if err != nil {
			logger.Errorf("decrypt username error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		password, err := dec.Decrypt(data.Password)
		if err != nil {
			logger.Errorf("decrypt password error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		notes, err := dec.Decrypt(data.Notes)
		if err != nil {
			logger.Errorf("decrypt notes error: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		resp := pb.PasswordReadResponse{
			Id:       data.ID,
			Name:     data.Name,
			Username: string(username),
			Password: string(password),
			Notes:    string(notes),
		}

		return &resp, nil
	}
}
