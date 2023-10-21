package note

import (
	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

func fromNoteWriteRequest(userID string, in *pb.NoteWriteRequest) storage.NoteData {
	return storage.NoteData{
		UserID: userID,
		Name:   in.Name,
		Notes:  in.Notes,
	}
}
