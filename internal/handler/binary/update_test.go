package binary

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

func TestGRPCUpdateHandler(t *testing.T) {

	binEncErr := errors.New("bin encrypt error")
	notesEncErr := errors.New("notes encrypt error")

	tests := []struct {
		name       string
		wantStatus codes.Code
		userErr    error
		ecnErr     error
		updErr     error
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
			name:       "update error",
			wantStatus: codes.Internal,
			updErr:     errors.New("update error"),
		},
		{
			name:       "already exists",
			wantStatus: codes.AlreadyExists,
			updErr:     storage.ErrWriteConflict,
		},
		{
			name:       "not found",
			wantStatus: codes.NotFound,
			updErr:     storage.ErrNoContent,
		},
		{
			name:       binEncErr.Error(),
			wantStatus: codes.Internal,
			ecnErr:     binEncErr,
		},
		{
			name:       notesEncErr.Error(),
			wantStatus: codes.Internal,
			ecnErr:     notesEncErr,
		},
	}

	for _, tcase := range tests {

		req := pb.BinaryUpdateRequest{
			Id: 1,
			Write: &pb.BinaryWriteRequest{
				Name:  "name",
				Size:  64,
				Notes: "notes",
			},
		}

		pu := BinaryUpdaterFunc(func(context.Context, storage.BinaryData) error {
			return tcase.updErr
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
			return nil, nil
		})

		t.Run(tcase.name, func(t *testing.T) {
			_, err := NewGRPCUpdateHandler(pu, getUserID, enc)(context.Background(), &req)
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
