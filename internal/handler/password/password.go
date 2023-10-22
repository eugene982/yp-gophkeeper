package password

import (
	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

func fromPasswordWriteRequest(userID string, id int32, in *pb.PasswordWriteRequest) storage.PasswordData {
	return storage.PasswordData{
		ID:       id,
		UserID:   userID,
		Name:     in.Name,
		Username: in.Username,
		Password: in.Password,
		Notes:    in.Notes,
	}
}

func toPasswordReadResponse(data storage.PasswordData) pb.PasswordReadResponse {
	return pb.PasswordReadResponse{
		Id:       data.ID,
		Name:     data.Name,
		Username: data.Username,
		Password: data.Password,
		Notes:    data.Notes,
	}
}
