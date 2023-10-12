package grpc

import (
	"context"
	"net"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/golang/protobuf/ptypes/empty"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler/ping"
	"github.com/eugene982/yp-gophkeeper/internal/handler/register"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
)

type ServerLogic interface {
	ping.Pinger
	register.Register
}

type GRPCServer struct {
	pb.UnimplementedGophKeeperServer

	listen net.Listener
	server *grpc.Server

	// handlers
	pingHandler ping.GRPCHahdler
	regHandler  register.GRPCHandler
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

	// создаём gRPC-сервер без зарегистрированной службы
	// с прослойками:
	//	- логирования
	//	- валидации входящих данных
	srv.server = grpc.NewServer(grpc.ChainUnaryInterceptor(
		loggerInterceptor,
		protovalidate_middleware.UnaryServerInterceptor(validator),
	))

	// Подключаем ручки
	srv.pingHandler = ping.NewRPCPingHandler(logic)
	srv.regHandler = register.NewRPCRegisterHandler(logic)

	// регистрируем сервис
	pb.RegisterGophKeeperServer(srv.server, srv)

	return &srv, nil
}

func (s *GRPCServer) Start() error {
	return s.server.Serve(s.listen)
}

// loggerInterceptor прослойка логирования запросов
func loggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()

	logger.Info("request",
		"method", info.FullMethod)

	logger.Debug("incoming",
		"request", req)

	resp, err = handler(ctx, req)

	logger.Info("request",
		"duration", time.Since(start),
		"ok", err == nil)

	logger.Debug("outgoing",
		"response", resp,
		"error", err)

	return
}

func (s GRPCServer) Ping(ctx context.Context, in *empty.Empty) (*pb.PingResponse, error) {
	if s.pingHandler == nil {
		return s.UnimplementedGophKeeperServer.Ping(ctx, in)
	}
	return s.pingHandler(ctx, in)
}

func (s GRPCServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if s.regHandler == nil {
		return s.UnimplementedGophKeeperServer.Register(ctx, in)
	}
	return s.regHandler(ctx, in)
}
