package grpc

import (
	"context"
	"net"

	"github.com/bufbuild/protovalidate-go"
	"github.com/golang/protobuf/ptypes/empty"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/handler/list"
	"github.com/eugene982/yp-gophkeeper/internal/handler/login"
	"github.com/eugene982/yp-gophkeeper/internal/handler/ping"
	"github.com/eugene982/yp-gophkeeper/internal/handler/register"
)

type ServerLogic interface {
	ping.Pinger
	register.Register
	list.ListGetter
}

// List implements list.List.

type GRPCServer struct {
	pb.UnimplementedGophKeeperServer

	listen net.Listener
	server *grpc.Server

	// handlers
	pingHandler  ping.GRPCHahdler
	regHandler   register.GRPCHandler
	loginHandler login.GRPCHandler
	listHandler  list.GRPCHandler
}

func NewServer(logic ServerLogic, addr, sekret_key string) (*GRPCServer, error) {
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
	srv.listHandler = list.NewRPCListHandler(logic, handler.NewMDUserIDGetter(sekret_key))

	// регистрируем сервис
	pb.RegisterGophKeeperServer(srv.server, srv)

	return &srv, nil
}

func (s *GRPCServer) Start() error {
	return s.server.Serve(s.listen)
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

func (s GRPCServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	if s.loginHandler == nil {
		return s.UnimplementedGophKeeperServer.Login(ctx, in)
	}
	return s.loginHandler(ctx, in)
}

func (s GRPCServer) List(ctx context.Context, in *empty.Empty) (*pb.ListResponse, error) {
	if s.listHandler == nil {
		return s.UnimplementedGophKeeperServer.List(ctx, in)
	}
	return s.listHandler(ctx, in)
}
