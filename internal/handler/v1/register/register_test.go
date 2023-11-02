package register

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

func TestRPCRegisterHandler(t *testing.T) {

	tests := []struct {
		name       string
		writeErr   error
		hashErr    error
		tokenErr   error
		wantStatus codes.Code
	}{
		{name: "ok", wantStatus: 0},
		{
			name:       "hash err",
			hashErr:    errors.New("hash err"),
			wantStatus: codes.Internal,
		},
		{
			name:       "token err",
			tokenErr:   errors.New("token err"),
			wantStatus: codes.Internal,
		},
		{
			name:       "already exists",
			writeErr:   storage.ErrWriteConflict,
			wantStatus: codes.AlreadyExists,
		},
		{
			name:       "write error",
			writeErr:   errors.New("write err"),
			wantStatus: codes.Internal,
		},
	}
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {

			w := UserWriterFunc(func(ctx context.Context, ud storage.UserData) error {
				return tcase.writeErr
			})

			hashFn := PasswordHashFunc(func(s string) (string, error) {
				if tcase.hashErr != nil {
					return "", tcase.hashErr
				}
				return "hash", nil
			})

			tokenFn := TokenGenFunc(func(s string) (string, error) {
				if tcase.tokenErr != nil {
					return "", tcase.tokenErr
				}
				return "token", nil
			})

			req := pb.RegisterRequest{
				Login:    "user",
				Password: "password",
			}

			resp, err := NewRPCRegisterHandler(w, hashFn, tokenFn)(context.Background(), &req)
			if tcase.wantStatus == 0 {
				assert.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, "token", resp.Token)
			} else {
				assert.Error(t, err)
			}

		})
	}
}
