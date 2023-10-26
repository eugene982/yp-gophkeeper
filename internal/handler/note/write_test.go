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

func TestGRPCWriteHandler(t *testing.T) {

	notesEncErr := errors.New("notes encrypt error")

	tests := []struct {
		name       string
		wantStatus codes.Code
		userErr    error
		ecnErr     error
		writeErr   error
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
			name:       "write error",
			wantStatus: codes.Internal,
			writeErr:   errors.New("write error"),
		},
		{
			name:       "already exists",
			wantStatus: codes.AlreadyExists,
			writeErr:   storage.ErrWriteConflict,
		},
		{
			name:       notesEncErr.Error(),
			wantStatus: codes.Internal,
			ecnErr:     notesEncErr,
		},
	}

	for _, tcase := range tests {

		req := pb.NoteWriteRequest{
			Name:  "name",
			Notes: "notes",
		}

		pw := NoteWritterFunc(func(context.Context, storage.NoteData) error {
			return tcase.writeErr
		})

		getUserID := handler.GetUserIDFunc(func(context.Context) (string, error) {
			if tcase.userErr != nil {
				return "", tcase.userErr
			}
			return "user", nil
		})

		enc := crypt.EncryptFunc(func(text []byte) ([]byte, error) {
			if tcase.ecnErr == notesEncErr && string(text) == "notes" {
				return nil, tcase.ecnErr
			}
			return text, nil
		})

		t.Run(tcase.name, func(t *testing.T) {
			_, err := NewGRPCWriteHandler(pw, getUserID, enc)(context.Background(), &req)
			if tcase.wantStatus == 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				status, ok := status.FromError(err)
				require.Equal(t, true, ok)
				assert.Equal(t, tcase.wantStatus, status.Code())
			}
		})
	}
}
