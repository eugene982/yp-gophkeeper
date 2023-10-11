package ping

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/stretchr/testify/assert"
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

			var pinger = handler.PingerFunc(func(context.Context) error {
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
