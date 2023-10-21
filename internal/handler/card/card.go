package card

import (
	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

func fromCardWriteRequest(userID string, in *pb.CardWriteRequest) storage.CardData {
	return storage.CardData{
		UserID: userID,
		Name:   in.Name,
		Number: in.Number,
		Notes:  in.Notes,
	}
}
