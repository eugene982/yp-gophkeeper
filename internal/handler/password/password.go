package password

import (
	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

func fromPasswordWriteRequest(userID string, in *pb.PasswordWriteRequest) storage.PasswordData {
	return storage.PasswordData{
		UserID:   userID,
		Name:     in.Name,
		Username: in.Username,
		Password: in.Password,
		Notes:    in.Notes,
	}
}

// func fromPasswordReadResponse(userID string, in *pb.PasswordReadResponse) storage.PasswordData {
// 	return storage.PasswordData{
// 		UserID:   userID,
// 		Name:     in.Name,
// 		Username: in.Username,
// 		Password: in.Password,
// 		Notes:    in.Notes,
// 	}
// }
