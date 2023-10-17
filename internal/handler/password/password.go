package password

import (
	"context"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PasswordManadger interface {
	PasswordLister
	PasswordWriter
	PasswordReader
	PasswordUpdater
	PasswordDeleter
}

type PasswordLister interface {
	PasswordList(context.Context, string) ([]string, error)
}

type PasswordWriter interface {
	PasswordWrite(context.Context, storage.PasswordData) error
}

type PasswordReader interface {
	PasswordRead(context.Context, string, string) (storage.PasswordData, error)
}

type PasswordUpdater interface {
	PasswordUpdate(context.Context, storage.PasswordData) error
}

type PasswordDeleter interface {
	PasswordDelete(context.Context, string, string) error
}

type GRPCListHandler func(context.Context, *empty.Empty) (*pb.PasswordListResponse, error)
type GRPCWriteHandler func(context.Context, *pb.PasswordWriteRequest) (*empty.Empty, error)
type GRPCUpdateHandler func(context.Context, *pb.PasswordWriteRequest) (*empty.Empty, error)
type GRPReadHandler func(context.Context, *pb.PasswordReadRequest) (*empty.Empty, error)
type GRPDeleteHandler func(context.Context, *pb.PasswordReadRequest) (*empty.Empty, error)

// NewRPCListHandler - ручка возвращает список наименований паролей
func NewRPCListHandler(ls PasswordLister, ug handler.UserIDGetter) GRPCListHandler {
	return func(ctx context.Context, e *empty.Empty) (*pb.PasswordListResponse, error) {
		var resp pb.PasswordListResponse

		userID, err := ug.GetUserID(ctx)
		if err != nil {
			return nil, err
		}

		resp.Names, err = ls.PasswordList(ctx, userID)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &resp, nil
	}
}
