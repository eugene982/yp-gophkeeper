package password

import (
	"context"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
	"github.com/golang/protobuf/ptypes/empty"
)

type PasswordWritter interface {
	PasswordWrite(ctx context.Context, data storage.PasswordData) error
}

type PasswordWritterFunc func(ctx context.Context, data any) error

func (f PasswordWritterFunc) PasswordWrite(ctx context.Context, data storage.PasswordData) error {
	return f(ctx, data)
}

var _ PasswordWritter = PasswordWritterFunc(nil)

type GRPCWriteHandler func(ctx context.Context, in *pb.PasswordRequest) (*empty.Empty, error)

// func NewGRPCWriteHandler(w PasswordWritter, ug handler.UserIDGetter) GRPCWriteHandler {
// 	return func(ctx context.Context, in *pb.PasswordRequest) (*empty.Empty, error) {
// 		userID, err := ug.GetUserID(ctx)
// 		if err != nil {
// 			return nil, err
// 		}

// 		err = w.PasswordWrite(ctx, toPasswordData(userID, in))
// 		if err != nil {
// 			return nil, status.Error(codes.Internal, err.Error())
// 		}

// 		return &empty.Empty{}, nil
// 	}
// }
