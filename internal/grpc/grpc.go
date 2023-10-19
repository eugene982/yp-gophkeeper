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
	"github.com/eugene982/yp-gophkeeper/internal/handler"
	"github.com/eugene982/yp-gophkeeper/internal/handler/card"
	"github.com/eugene982/yp-gophkeeper/internal/handler/list"
	"github.com/eugene982/yp-gophkeeper/internal/handler/login"
	"github.com/eugene982/yp-gophkeeper/internal/handler/note"
	"github.com/eugene982/yp-gophkeeper/internal/handler/password"
	"github.com/eugene982/yp-gophkeeper/internal/handler/ping"
	"github.com/eugene982/yp-gophkeeper/internal/handler/register"
	"github.com/eugene982/yp-gophkeeper/internal/storage"
)

var (
	TOKEN_SECRET_KEY = "sekret=key"
	TOKEN_EXP        = time.Hour
	PASSWORD_SALT    = "password=salt"
)

// type ServerLogic interface {
// 	ping.Pinger
// 	register.Register
// 	list.ListGetter
// 	login.Login

// 	password.PasswordStorage
// }

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

	// password
	paswordListHandler password.GRPCListHandler

	// cards
	cardListHandler card.GRPCListHandler

	// notes
	noteListHandler note.GRPCListHandler
}

func NewServer(store storage.Storage, addr string) (*GRPCServer, error) {
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

	// Функция хеширования паролей
	hashFn := func(passwd string) (string, error) {
		return handler.PasswordHash(passwd, PASSWORD_SALT)
	}
	// Функция генерирования токена
	tokenFn := func(userId string) (string, error) {
		return handler.MakeToken(userId, TOKEN_SECRET_KEY, TOKEN_EXP)
	}
	// Функция сравнения хеша и пароля пользователя
	checkFn := func(password, hash string) bool {
		return handler.CheckPasswordHash(hash, password, PASSWORD_SALT)
	}

	getUserID := func(ctx context.Context) (string, error) {
		return handler.GetUserIDFromMD(ctx, TOKEN_SECRET_KEY)
	}

	// Подключаем ручки
	srv.pingHandler = ping.NewRPCPingHandler(store)
	srv.regHandler = register.NewRPCRegisterHandler(store, hashFn, tokenFn)
	srv.loginHandler = login.NewRPCLoginHandler(store, checkFn, tokenFn)
	srv.listHandler = list.NewRPCListHandler(store, getUserID)

	// Password
	srv.paswordListHandler = password.NewGRPCListHandler(store, getUserID)

	// Payment card
	srv.cardListHandler = card.NewGRPCListHandler(store, getUserID)

	// Notes
	srv.noteListHandler = note.NewGRPCListHandler(store, getUserID)

	// регистрируем сервис
	pb.RegisterGophKeeperServer(srv.server, srv)

	return &srv, nil
}

func (s *GRPCServer) Start() error {
	return s.server.Serve(s.listen)
}

func (s GRPCServer) Ping(ctx context.Context, in *empty.Empty) (*pb.PingResponse, error) {
	if s.pingHandler != nil {
		return s.pingHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.Ping(ctx, in)
}

func (s GRPCServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if s.regHandler != nil {
		return s.regHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.Register(ctx, in)
}

func (s GRPCServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	if s.loginHandler != nil {
		return s.loginHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.Login(ctx, in)
}

func (s GRPCServer) List(ctx context.Context, in *empty.Empty) (*pb.ListResponse, error) {
	if s.listHandler != nil {
		return s.listHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.List(ctx, in)
}

// Password

func (s GRPCServer) PasswordList(ctx context.Context, in *empty.Empty) (*pb.PasswordListResponse, error) {
	if s.paswordListHandler != nil {
		return s.paswordListHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.PasswordList(ctx, in)
}

func (s GRPCServer) PasswordWrite(ctx context.Context, in *pb.PasswordRequest) (*empty.Empty, error) {
	return s.UnimplementedGophKeeperServer.PasswordWrite(ctx, in)
}

func (s GRPCServer) PasswordRead(ctx context.Context, in *pb.PasswordReadRequest) (*pb.PasswordResponse, error) {
	return s.UnimplementedGophKeeperServer.PasswordRead(ctx, in)
}

func (s GRPCServer) PasswordUpdate(ctx context.Context, in *pb.PasswordRequest) (*empty.Empty, error) {
	return s.UnimplementedGophKeeperServer.PasswordUpdate(ctx, in)
}

func (s GRPCServer) PasswordDelete(ctx context.Context, in *pb.PasswordReadRequest) (*empty.Empty, error) {
	return s.UnimplementedGophKeeperServer.PasswordDelete(ctx, in)
}

// Cards

func (s GRPCServer) CardList(ctx context.Context, in *empty.Empty) (*pb.CardListResponse, error) {
	if s.cardListHandler != nil {
		return s.cardListHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.CardList(ctx, in)
}

// Notes

func (s GRPCServer) NoteList(ctx context.Context, in *empty.Empty) (*pb.NoteListResponse, error) {
	if s.noteListHandler != nil {
		return s.noteListHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.NoteList(ctx, in)
}
