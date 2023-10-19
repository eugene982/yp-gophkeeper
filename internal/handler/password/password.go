package password

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

type GRPCManadger interface {
	PasswordList(ctx context.Context, in *empty.Empty) (*pb.PasswordListResponse, error)
	PasswordWrite(ctx context.Context, in *pb.PasswordRequest) (*empty.Empty, error)
	PasswordRead(ctx context.Context, in *pb.PasswordReadRequest) (*pb.PasswordResponse, error)
	PasswordUpdate(ctx context.Context, in *pb.PasswordRequest) (*empty.Empty, error)
}

type PasswordStorage interface {
	//NamesList(ctx context.Context, tab storage.TableName, userID string) ([]string, error)
	//Write(ctx context.Context, data any) error
	Update(ctx context.Context, data any) error
	ReadByName(ctx context.Context, tab storage.TableName, userID, name string) (any, error)
	DeleteByName(ctx context.Context, tab storage.TableName, userID, name string) error
}

//func NewPasswordStore

type passwordManadger struct {
	tabname storage.TableName
	store   PasswordStorage
	//userGetter handler.UserIDGetter
}

// NewRPCListHandler - ручка возвращает список наименований паролей
// func NewGRPCPasswordManadger(s PasswordStorage, ug handler.UserIDGetter) GRPCManadger {
// 	return passwordManadger{
// 		tabname:    "passwords",
// 		store:      s,
// 		userGetter: ug,
// 	}
// }

// PasswordWrite запись нового пароля
// func (m passwordManadger) PasswordWrite(ctx context.Context, in *pb.PasswordRequest) (*empty.Empty, error) {
// 	userID, err := m.userGetter.GetUserID(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = m.store.Write(ctx, toPasswordData(userID, in))
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &empty.Empty{}, nil
// }

// func (m passwordManadger) PasswordRead(ctx context.Context, in *pb.PasswordReadRequest) (*pb.PasswordResponse, error) {
// 	userID, err := m.userGetter.GetUserID(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	data, err := m.store.ReadByName(ctx, m.tabname, userID, in.Name)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	} else if res, ok := data.(storage.PasswordData); ok {
// 		return &pb.PasswordResponse{
// 			Name:     res.Name,
// 			Username: res.Username,
// 			Password: res.Password,
// 			Notes:    res.Notes,
// 		}, nil
// 	}
// 	return nil, status.Error(codes.Internal, "unknown data type")
// }

// // PasswordUpdate обновление пользовательского пароля
// func (m passwordManadger) PasswordUpdate(ctx context.Context, in *pb.PasswordRequest) (*empty.Empty, error) {
// 	userID, err := m.userGetter.GetUserID(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = m.store.Update(ctx, toPasswordData(userID, in)); err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
// 	return &empty.Empty{}, nil
// }

// func toPasswordData(userID string, in *pb.PasswordRequest) storage.PasswordData {
// 	return storage.PasswordData{
// 		UserID:   userID,
// 		Name:     in.Name,
// 		Username: in.Username,
// 		Password: in.Password,
// 		Notes:    in.Notes,
// 	}
// }
