package password

import (
	"context"
	"errors"
	"testing"

	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCListHandler(t *testing.T) {

	tests := []struct {
		name       string
		wantStatus codes.Code
		userErr    error
		listErr    error
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
			name:       "list error",
			wantStatus: codes.Internal,
			listErr:    errors.New("list error"),
		},
		{
			name:       "not found",
			wantStatus: codes.NotFound,
			listErr:    storage.ErrNoContent,
		},
	}

	for _, tcase := range tests {

		list := []string{"p1", "p2"}

		getList := PasswordListGetterFunc(func(ctx context.Context, userID string) ([]string, error) {
			if tcase.listErr != nil {
				return nil, tcase.listErr
			}
			return list, nil
		})

		getUserID := handler.GetUserIDFunc(func(context.Context) (string, error) {
			if tcase.userErr != nil {
				return "", tcase.userErr
			}
			return "user", nil
		})

		resp, err := NewGRPCListHandler(getList, getUserID)(context.Background(), nil)

		t.Run(tcase.name, func(t *testing.T) {
			if tcase.wantStatus == 0 {
				require.NoError(t, err)
				assert.Equal(t, list, resp.Names)

			} else {
				assert.Error(t, err)
				status, ok := status.FromError(err)
				require.Equal(t, true, ok)
				assert.Equal(t, tcase.wantStatus, status.Code())
			}
		})
	}
}
