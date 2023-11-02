package login

import (
	"context"
	"errors"
	"testing"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestRPCRegisterHandler(t *testing.T) {

	tests := []struct {
		name       string
		login      string
		wantStatus codes.Code
	}{
		{name: "ok", login: "user", wantStatus: 0},
		{name: "unauthenticated", login: "user", wantStatus: codes.Unauthenticated},
		{name: "read user error", login: "user", wantStatus: codes.Internal},
		{name: "no content", login: "user", wantStatus: codes.Unauthenticated},
		{name: "token error", login: "user", wantStatus: codes.Internal},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {

			reader := UserReaderFunc(func(ctx context.Context, s string) (res storage.UserData, err error) {
				if tcase.name == "read user error" {
					err = errors.New(tcase.name)
				} else if tcase.name == "no content" {
					err = storage.ErrNoContent
				} else {
					res.UserID = tcase.login
					res.PasswordHash = "hash"
				}
				return
			})

			checkFn := HashCheckFunc(func(s1, s2 string) bool {
				return tcase.wantStatus != codes.Unauthenticated
			})

			req := pb.LoginRequest{
				Login:    tcase.login,
				Password: tcase.login,
			}

			token := TokenGenFunc(func(s string) (string, error) {
				if tcase.name == "token error" {
					return "", errors.New(tcase.name)
				}
				return "token", nil
			})

			resp, err := NewRPCLoginHandler(reader, checkFn, token)(context.Background(), &req)
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
