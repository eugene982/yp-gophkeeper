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
	"github.com/eugene982/yp-gophkeeper/internal/handler/binary"
	"github.com/eugene982/yp-gophkeeper/internal/handler/card"
	"github.com/eugene982/yp-gophkeeper/internal/handler/list"
	"github.com/eugene982/yp-gophkeeper/internal/handler/login"
	"github.com/eugene982/yp-gophkeeper/internal/handler/note"
	"github.com/eugene982/yp-gophkeeper/internal/handler/password"
	"github.com/eugene982/yp-gophkeeper/internal/handler/ping"
	"github.com/eugene982/yp-gophkeeper/internal/handler/register"
	"github.com/eugene982/yp-gophkeeper/internal/storage"

	crypt "github.com/eugene982/yp-gophkeeper/internal/crypto"
)

var (
	TOKEN_SECRET_KEY = "sekret=key"
	TOKEN_EXP        = time.Hour
	PASSWORD_SALT    = "password=salt"
	CRYPTO_KEY       = []byte("GopherKeeperKey!") // 16, 24, 34 байта
)

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
	passwdListHandler   password.GRPCListHandler
	passwdWriteHandler  password.GRPCWriteHandler
	passwdReadHandler   password.GRPCReadHandler
	passwdDeleteHandler password.GRPCDeleteHandler
	passwdUpdateHandler password.GRPCUpdateHandler

	// cards
	cardListHandler   card.GRPCListHandler
	cardWriteHandler  card.GRPCWriteHandler
	cardReadHandler   card.GRPCReadHandler
	cardDeleteHandler card.GRPCDeleteHandler
	cardUpdateHandler card.GRPCUpdateHandler

	// notes
	noteListHandler   note.GRPCListHandler
	noteWriteHandler  note.GRPCWriteHandler
	noteReadHandler   note.GRPCReadHandler
	noteDeleteHandler note.GRPCDeleteHandler
	noteUpdateHandler note.GRPCUpdateHandler

	// binary
	binaryListHandler     binary.GRPCListHandler
	binaryWriteHandler    binary.GRPCWriteHandler
	binaryReadHandler     binary.GRPCReadHandler
	binaryDeleteHandler   binary.GRPCDeleteHandler
	binaryUpdateHandler   binary.GRPCUpdateHandler
	binaryUploadHandler   binary.GRPCUploadHandler
	binaryDownloadHandler binary.GRPCDownloadHandler
}

// NewServer функция-коструктор нового grps сервера
func NewServer(store storage.Storage, crypt crypt.EncryptDecryptor, addr string) (*GRPCServer, error) {
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
	srv.server = grpc.NewServer(
		grpc.ChainStreamInterceptor(
			loggerStreamInterceptor,
			protovalidate_middleware.StreamServerInterceptor(validator)),
		grpc.ChainUnaryInterceptor(
			loggerInterceptor,
			protovalidate_middleware.UnaryServerInterceptor(validator)),
	)

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

	// Функция вытаскивания ид. пользователя из контекста
	getUserID := func(ctx context.Context) (string, error) {
		return handler.GetUserIDFromMD(ctx, TOKEN_SECRET_KEY)
	}

	// Подключаем ручки
	srv.pingHandler = ping.NewRPCPingHandler(store)
	srv.regHandler = register.NewRPCRegisterHandler(store, hashFn, tokenFn)
	srv.loginHandler = login.NewRPCLoginHandler(store, checkFn, tokenFn)
	srv.listHandler = list.NewRPCListHandler(store, getUserID)

	// Password
	srv.passwdListHandler = password.NewGRPCListHandler(store, getUserID)
	srv.passwdWriteHandler = password.NewGRPCWriteHandler(store, getUserID, crypt)
	srv.passwdReadHandler = password.NewGRPCReadHandler(store, getUserID, crypt)
	srv.passwdDeleteHandler = password.NewGRPCDeleteHandler(store, getUserID)
	srv.passwdUpdateHandler = password.NewGRPCUpdateHandler(store, getUserID, crypt)

	// Payment card
	srv.cardListHandler = card.NewGRPCListHandler(store, getUserID)
	srv.cardWriteHandler = card.NewGRPCWriteHandler(store, getUserID, crypt)
	srv.cardReadHandler = card.NewGRPCReadHandler(store, getUserID, crypt)
	srv.cardDeleteHandler = card.NewGRPCDeleteHandler(store, getUserID)
	srv.cardUpdateHandler = card.NewGRPCUpdateHandler(store, getUserID, crypt)

	// Notes
	srv.noteListHandler = note.NewGRPCListHandler(store, getUserID)
	srv.noteWriteHandler = note.NewGRPCWriteHandler(store, getUserID, crypt)
	srv.noteReadHandler = note.NewGRPCReadHandler(store, getUserID, crypt)
	srv.noteDeleteHandler = note.NewGRPCDeleteHandler(store, getUserID)
	srv.noteUpdateHandler = note.NewGRPCUpdateHandler(store, getUserID, crypt)

	// binary
	srv.binaryListHandler = binary.NewGRPCListHandler(store, getUserID)
	srv.binaryWriteHandler = binary.NewGRPCWriteHandler(store, getUserID, crypt)
	srv.binaryReadHandler = binary.NewGRPCReadHandler(store, getUserID, crypt)
	srv.binaryDeleteHandler = binary.NewGRPCDeleteHandler(store, getUserID)
	srv.binaryUpdateHandler = binary.NewGRPCUpdateHandler(store, getUserID, crypt)
	srv.binaryUploadHandler = binary.NewGRPCUploaderHandler(store)
	srv.binaryDownloadHandler = binary.NewGRPCDownloadHandler(store)

	// регистрируем сервис
	pb.RegisterGophKeeperServer(srv.server, srv)

	return &srv, nil
}

// Start - запуск сервера
func (s *GRPCServer) Start() error {
	return s.server.Serve(s.listen)
}

// Stop - остановка сервера
func (s *GRPCServer) Stop() {
	s.server.Stop()
}

func (s GRPCServer) Ping(ctx context.Context, in *empty.Empty) (*pb.PingResponse, error) {
	if s.pingHandler != nil {
		return s.pingHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.Ping(ctx, in)
}

// User

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
	if s.passwdListHandler != nil {
		return s.passwdListHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.PasswordList(ctx, in)
}

func (s GRPCServer) PasswordWrite(ctx context.Context, in *pb.PasswordWriteRequest) (*empty.Empty, error) {
	if s.passwdWriteHandler != nil {
		return s.passwdWriteHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.PasswordWrite(ctx, in)
}

func (s GRPCServer) PasswordRead(ctx context.Context, in *pb.PasswordReadRequest) (*pb.PasswordReadResponse, error) {
	if s.passwdReadHandler != nil {
		return s.passwdReadHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.PasswordRead(ctx, in)
}

func (s GRPCServer) PasswordUpdate(ctx context.Context, in *pb.PasswordUpdateRequest) (*empty.Empty, error) {
	if s.passwdUpdateHandler != nil {
		return s.passwdUpdateHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.PasswordUpdate(ctx, in)
}

func (s GRPCServer) PasswordDelete(ctx context.Context, in *pb.PasswordDelRequest) (*empty.Empty, error) {
	if s.passwdDeleteHandler != nil {
		return s.passwdDeleteHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.PasswordDelete(ctx, in)
}

// Cards

func (s GRPCServer) CardList(ctx context.Context, in *empty.Empty) (*pb.CardListResponse, error) {
	if s.cardListHandler != nil {
		return s.cardListHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.CardList(ctx, in)
}

func (s GRPCServer) CardWrite(ctx context.Context, in *pb.CardWriteRequest) (*empty.Empty, error) {
	if s.cardWriteHandler != nil {
		return s.cardWriteHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.CardWrite(ctx, in)
}

func (s GRPCServer) CardRead(ctx context.Context, in *pb.CardReadRequest) (*pb.CardReadResponse, error) {
	if s.cardReadHandler != nil {
		return s.cardReadHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.CardRead(ctx, in)
}

func (s GRPCServer) CardUpdate(ctx context.Context, in *pb.CardUpdateRequest) (*empty.Empty, error) {
	if s.cardUpdateHandler != nil {
		return s.cardUpdateHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.CardUpdate(ctx, in)
}

func (s GRPCServer) CardDelete(ctx context.Context, in *pb.CardDelRequest) (*empty.Empty, error) {
	if s.cardDeleteHandler != nil {
		return s.cardDeleteHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.CardDelete(ctx, in)
}

// Notes

func (s GRPCServer) NoteList(ctx context.Context, in *empty.Empty) (*pb.NoteListResponse, error) {
	if s.noteListHandler != nil {
		return s.noteListHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.NoteList(ctx, in)
}

func (s GRPCServer) NoteWrite(ctx context.Context, in *pb.NoteWriteRequest) (*empty.Empty, error) {
	if s.noteWriteHandler != nil {
		return s.noteWriteHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.NoteWrite(ctx, in)
}

func (s GRPCServer) NoteRead(ctx context.Context, in *pb.NoteReadRequest) (*pb.NoteReadResponse, error) {
	if s.noteReadHandler != nil {
		return s.noteReadHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.NoteRead(ctx, in)
}

func (s GRPCServer) NoteUpdate(ctx context.Context, in *pb.NoteUpdateRequest) (*empty.Empty, error) {
	if s.noteUpdateHandler != nil {
		return s.noteUpdateHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.NoteUpdate(ctx, in)
}

func (s GRPCServer) NoteDelete(ctx context.Context, in *pb.NoteDelRequest) (*empty.Empty, error) {
	if s.noteDeleteHandler != nil {
		return s.noteDeleteHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.NoteDelete(ctx, in)
}

// Binaries

func (s GRPCServer) BinaryList(ctx context.Context, in *empty.Empty) (*pb.BinaryListResponse, error) {
	if s.binaryListHandler != nil {
		return s.binaryListHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.BinaryList(ctx, in)
}

func (s GRPCServer) BinaryWrite(ctx context.Context, in *pb.BinaryWriteRequest) (*pb.BinaryWriteResponse, error) {
	if s.binaryWriteHandler != nil {
		return s.binaryWriteHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.BinaryWrite(ctx, in)
}

func (s GRPCServer) BinaryRead(ctx context.Context, in *pb.BinaryReadRequest) (*pb.BinaryReadResponse, error) {
	if s.binaryReadHandler != nil {
		return s.binaryReadHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.BinaryRead(ctx, in)
}

func (s GRPCServer) BinaryUpdate(ctx context.Context, in *pb.BinaryUpdateRequest) (*empty.Empty, error) {
	if s.binaryUpdateHandler != nil {
		return s.binaryUpdateHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.BinaryUpdate(ctx, in)
}

func (s GRPCServer) BinaryDelete(ctx context.Context, in *pb.BinaryDelRequest) (*empty.Empty, error) {
	if s.binaryDeleteHandler != nil {
		return s.binaryDeleteHandler(ctx, in)
	}
	return s.UnimplementedGophKeeperServer.BinaryDelete(ctx, in)
}

func (s GRPCServer) BinaryUpload(us pb.GophKeeper_BinaryUploadServer) error {
	if s.binaryUploadHandler != nil {
		return s.binaryUploadHandler(us)
	}
	return s.UnimplementedGophKeeperServer.BinaryUpload(us)
}

func (s GRPCServer) BinaryDownload(req *pb.BidaryDownloadRequest, ds pb.GophKeeper_BinaryDownloadServer) error {
	if s.binaryDownloadHandler != nil {
		return s.binaryDownloadHandler(req, ds)
	}
	return s.UnimplementedGophKeeperServer.BinaryDownload(req, ds)
}
