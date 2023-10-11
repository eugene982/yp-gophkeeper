package grpc

import (
	"context"
	"net"

	"github.com/bufbuild/protovalidate-go"
	"github.com/golang/protobuf/ptypes/empty"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler/ping"
)

type ServerLogic interface {
	ping.Pinger
}

type GRPCServer struct {
	pb.GophKeeperServer

	listen net.Listener
	server *grpc.Server

	// handlers
	pingHandler ping.GRPCHahdler
}

func NewServer(logic ServerLogic, addr string) (*GRPCServer, error) {
	var (
		srv GRPCServer
		err error
	)

	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	// определяем адрес сервера
	srv.listen, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	// создаём gRPC-сервер без зарегистрированной службы с прослойкой валидации входящих данных
	srv.server = grpc.NewServer(grpc.UnaryInterceptor(
		protovalidate_middleware.UnaryServerInterceptor(validator),
	))

	//srv.pingHandler = ping.NewRPCPingHandler()

	// регистрируем сервис
	pb.RegisterGophKeeperServer(srv.server, srv)

	return &srv, nil
}

func (s *GRPCServer) Start() error {
	return s.server.Serve(s.listen)
}

func (s GRPCServer) Ping(ctx context.Context, in *empty.Empty) (*pb.PingResponse, error) {
	if s.pingHandler == nil {
		return s.GophKeeperServer.Ping(ctx, in)
	}
	return s.pingHandler(ctx, in)
}
