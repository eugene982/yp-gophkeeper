package ping

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestPing(t *testing.T) {

	tests := []struct {
		name       string
		wantStatus int
	}{
		{name: "ok", wantStatus: 200},
		{name: "internal error", wantStatus: 500},
	}
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/ping", nil)

			var pinger = PingerFunc(func(context.Context) error {
				if tcase.wantStatus == 200 {
					return nil
				} else {
					return fmt.Errorf("mock ping error")
				}
			})

			NewPingHandler(pinger).ServeHTTP(w, r)
			assert.Equal(t, tcase.wantStatus, w.Code)
		})
	}
}

func TestRPCPing(t *testing.T) {

	tests := []struct {
		name       string
		wantStatus codes.Code
	}{
		{name: "ok", wantStatus: 0},
		{name: "internal error", wantStatus: codes.Internal},
	}
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {

			var pinger = PingerFunc(func(context.Context) error {
				if tcase.wantStatus == 0 {
					return nil
				} else {
					return fmt.Errorf("error")
				}
			})

			resp, err := NewRPCPingHandler(pinger)(context.Background(), &empty.Empty{})
			if tcase.wantStatus == 0 {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			} else {
				assert.Error(t, err)
			}

		})
	}
}
