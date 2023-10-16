package list

import (
	"context"
	"errors"
	"testing"

	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRPCListHandler(t *testing.T) {

	tests := []struct {
		name       string
		data       storage.ListData
		wantStatus codes.Code
	}{
		{name: "ok", wantStatus: 0},
		{name: "internal error", wantStatus: codes.Internal},
		{name: "unauthenticated", wantStatus: codes.Unauthenticated},
		{name: "invalid argument", wantStatus: codes.InvalidArgument},
	}
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {

			list := ListGetterFunc(func(ctx context.Context, userID string) (storage.ListData, error) {
				if tcase.wantStatus == codes.Internal {
					return tcase.data, errors.New("some error")
				}
				return tcase.data, nil
			})

			getter := handler.UserIDGetterFunc(func(ctx context.Context) (string, error) {
				switch tcase.wantStatus {
				case codes.Unauthenticated:
					return "", status.Errorf(codes.Unauthenticated, "Unauthenticated")
				case codes.InvalidArgument:
					return "", status.Errorf(codes.InvalidArgument, "invalid argument")
				default:
					return "user", nil
				}
			})

			resp, err := NewRPCListHandler(list, getter)(context.Background(), nil)
			if tcase.wantStatus == 0 {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			} else {
				assert.Error(t, err)
			}

		})
	}
}
