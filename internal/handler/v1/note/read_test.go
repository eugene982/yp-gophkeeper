package note

import (
	"context"
	"errors"
	"testing"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	crypt "github.com/eugene982/yp-gophkeeper/internal/crypto"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCreadHandler(t *testing.T) {

	notesDecErr := errors.New("notes decrypt error")

	tests := []struct {
		name       string
		wantStatus codes.Code
		userErr    error
		decErr     error
		readErr    error
	}{
		{
			name: "ok",
		},
		{
			name:       "unauthenticated",
			wantStatus: codes.Unauthenticated,
			userErr:    handler.ErrRPCInvalidToken,
		},
		{
			name:       "read error",
			wantStatus: codes.Internal,
			readErr:    errors.New("read error"),
		},
		{
			name:       "not found",
			wantStatus: codes.NotFound,
			readErr:    storage.ErrNoContent,
		},
		{
			name:       notesDecErr.Error(),
			wantStatus: codes.Internal,
			decErr:     notesDecErr,
		},
	}

	for _, tcase := range tests {

		req := pb.NoteReadRequest{
			Name: "name",
		}

		pr := NoteReaderFunc(func(context.Context, string, string) (res storage.NoteData, err error) {
			err = tcase.readErr
			if err == nil {
				res.ID = 1
				res.UserID = "user_id"
				res.Name = "name"
				res.Notes = []byte("notes")
			}
			return
		})

		getUserID := handler.GetUserIDFunc(func(context.Context) (string, error) {
			if tcase.userErr != nil {
				return "", tcase.userErr
			}
			return "user", nil
		})

		dec := crypt.DecryptFunc(func(text []byte) ([]byte, error) {
			if tcase.decErr == notesDecErr && string(text) == "notes" {
				return nil, tcase.decErr
			}
			return text, nil
		})

		resp, err := NewGRPCReadHandler(pr, getUserID, dec)(context.Background(), &req)

		t.Run(tcase.name, func(t *testing.T) {
			if tcase.wantStatus == 0 {
				assert.NoError(t, err)

				assert.Equal(t, int64(1), resp.Id)
				assert.Equal(t, "name", resp.Name)
				assert.Equal(t, "notes", resp.Notes)

			} else {
				assert.Error(t, err)
				status, ok := status.FromError(err)
				require.Equal(t, true, ok)
				assert.Equal(t, tcase.wantStatus, status.Code())
			}
		})

	}

}
