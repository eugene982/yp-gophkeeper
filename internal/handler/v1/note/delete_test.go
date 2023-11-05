package note

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

func TestGRPCDeleteHandler(t *testing.T) {

	tests := []struct {
		name       string
		wantStatus codes.Code
		userErr    error
		delErr     error
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
			name:       "delete error",
			wantStatus: codes.Internal,
			delErr:     errors.New("delete error"),
		},
		{
			name:       "not found",
			wantStatus: codes.NotFound,
			delErr:     storage.ErrNoContent,
		},
	}

	for _, tcase := range tests {
		delete := NoteDeleteFunc(func(ctx context.Context, userID, name string) error {
			return tcase.delErr
		})

		getUserID := handler.GetUserIDFunc(func(context.Context) (string, error) {
			if tcase.userErr != nil {
				return "", tcase.userErr
			}
			return "user", nil
		})

		req := pb.NoteDelRequest{
			Name: "note",
		}
		_, err := NewGRPCDeleteHandler(delete, getUserID)(context.Background(), &req)

		t.Run(tcase.name, func(t *testing.T) {
			if tcase.wantStatus == 0 {
				require.NoError(t, err)

			} else {
				assert.Error(t, err)
				status, ok := status.FromError(err)
				require.Equal(t, true, ok)
				assert.Equal(t, tcase.wantStatus, status.Code())
			}
		})
	}

}
