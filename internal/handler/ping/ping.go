package ping

import (
	"context"
	"net/http"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
)

// Pinger интерфейс проверки соединения
type Pinger interface {
	Ping(context.Context) error
}

// Тип
type PingerFunc func(context.Context) error

func (f PingerFunc) Ping(ctx context.Context) error {
	return f(ctx)
}

var _ Pinger = PingerFunc(nil)

func NewPingHandler(pinger Pinger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		if err := pinger.Ping(r.Context()); err != nil {
			logger.Errorf("error ping handler: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))

	}
	return http.HandlerFunc(fn)
}

type GRPCHahdler func(context.Context, *empty.Empty) (*pb.PingResponse, error)

func NewRPCPingHandler(pinger Pinger) GRPCHahdler {

	pingResp := pb.PingResponse{
		Message: "pong",
	}

	fn := func(ctx context.Context, _ *empty.Empty) (*pb.PingResponse, error) {

		if err := pinger.Ping(ctx); err != nil {
			logger.Errorf("error ping rpc handler: %w", err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &pingResp, nil
	}

	return GRPCHahdler(fn)
}
