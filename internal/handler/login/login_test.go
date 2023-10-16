package login

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
)

func TestRPCRegisterHandler(t *testing.T) {

	tests := []struct {
		name       string
		login      string
		wantStatus codes.Code
	}{
		{name: "ok", login: "user", wantStatus: 0},
		{name: "unauthenticated", login: "user", wantStatus: codes.Unauthenticated},
		{name: "internal error", login: "user", wantStatus: codes.Internal},
	}
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {

			login := LoginFunc(func(ctx context.Context, login, _ string) (string, error) {
				switch tcase.wantStatus {
				case 0:
					return login, nil
				case codes.Unauthenticated:
					return "", handler.ErrUnauthenticated
				default:
					return "", fmt.Errorf("error")
				}
			})

			req := pb.LoginRequest{
				Login:    tcase.login,
				Password: tcase.login,
			}

			resp, err := NewRPCLoginHandler(login)(context.Background(), &req)
			if tcase.wantStatus == 0 {
				assert.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, tcase.login, resp.Token)
			} else {
				assert.Error(t, err)
			}

		})
	}
}
