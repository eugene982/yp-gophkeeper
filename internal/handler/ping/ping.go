package ping

import (
	"net/http"

	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
)

func NewPingHandler(pinger handler.Pinger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := pinger.Ping(r.Context())
		if err != nil {
			logger.Errorf("error ping handler: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}
	return http.HandlerFunc(fn)
}
