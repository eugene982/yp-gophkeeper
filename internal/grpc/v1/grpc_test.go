package grpc

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/handler/v1/binary"
	"github.com/eugene982/yp-gophkeeper/internal/handler/v1/card"
	"github.com/eugene982/yp-gophkeeper/internal/handler/v1/list"
	"github.com/eugene982/yp-gophkeeper/internal/handler/v1/login"
	"github.com/eugene982/yp-gophkeeper/internal/handler/v1/note"
	"github.com/eugene982/yp-gophkeeper/internal/handler/v1/password"
	"github.com/eugene982/yp-gophkeeper/internal/handler/v1/ping"
	"github.com/eugene982/yp-gophkeeper/internal/handler/v1/register"
)

func TestNewGRPCServer(t *testing.T) {

	server, err := NewServer(nil, nil, ":8080")
	require.NoError(t, err)
	require.NotNil(t, server)
}

func TestServerHandlers(t *testing.T) {

	server := GRPCServer{}
	ctx := context.Background()
	emt := &empty.Empty{}

	t.Run("ping", func(t *testing.T) {
		_, err := server.Ping(ctx, emt)
		require.Error(t, err)

		server.pingHandler = ping.GRPCHahdler(func(ctx context.Context, e *empty.Empty) (*pb.PingResponse, error) {
			return nil, status.Error(codes.Internal, "ping error")
		})

		_, err = server.Ping(ctx, emt)
		require.Error(t, err)
	})

	t.Run("register", func(t *testing.T) {
		_, err := server.Register(ctx, nil)
		require.Error(t, err)

		server.regHandler = register.GRPCHandler(func(context.Context, *pb.RegisterRequest) (*pb.RegisterResponse, error) {
			return nil, status.Error(codes.Internal, "reg error")
		})

		_, err = server.Register(ctx, nil)
		require.Error(t, err)
	})

	t.Run("login", func(t *testing.T) {
		_, err := server.Login(ctx, nil)
		require.Error(t, err)

		server.loginHandler = login.GRPCHandler(func(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error) {
			return nil, status.Error(codes.Internal, "login error")
		})

		_, err = server.Login(ctx, nil)
		require.Error(t, err)
	})

	t.Run("list", func(t *testing.T) {
		_, err := server.List(ctx, nil)
		require.Error(t, err)

		server.listHandler = list.GRPCHandler(func(context.Context, *empty.Empty) (*pb.ListResponse, error) {
			return nil, status.Error(codes.Internal, "list error")
		})

		_, err = server.List(ctx, nil)
		require.Error(t, err)
	})

	// password

	t.Run("password list", func(t *testing.T) {
		_, err := server.PasswordList(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "password list error")
		server.passwdListHandler = password.GRPCListHandler(func(ctx context.Context, in *empty.Empty) (*pb.PasswordListResponse, error) {
			return nil, resperr
		})

		_, err = server.PasswordList(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("password write", func(t *testing.T) {
		_, err := server.PasswordWrite(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "password write error")
		server.passwdWriteHandler = password.GRPCWriteHandler(func(ctx context.Context, in *pb.PasswordWriteRequest) (*empty.Empty, error) {
			return nil, resperr
		})

		_, err = server.PasswordWrite(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("password read", func(t *testing.T) {
		_, err := server.PasswordRead(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "password read error")
		server.passwdReadHandler = password.GRPCReadHandler(func(ctx context.Context, prr *pb.PasswordReadRequest) (*pb.PasswordReadResponse, error) {
			return nil, resperr
		})

		_, err = server.PasswordRead(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("password update", func(t *testing.T) {
		_, err := server.PasswordUpdate(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "password update error")
		server.passwdUpdateHandler = password.GRPCUpdateHandler(func(ctx context.Context, pur *pb.PasswordUpdateRequest) (*empty.Empty, error) {
			return nil, resperr
		})

		_, err = server.PasswordUpdate(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("password delete", func(t *testing.T) {
		_, err := server.PasswordDelete(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "password delete error")
		server.passwdDeleteHandler = password.GRPCDeleteHandler(func(ctx context.Context, in *pb.PasswordDelRequest) (*empty.Empty, error) {
			return nil, resperr
		})

		_, err = server.PasswordDelete(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	// card

	t.Run("card list", func(t *testing.T) {
		_, err := server.CardList(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "card list error")
		server.cardListHandler = card.GRPCListHandler(func(ctx context.Context, in *empty.Empty) (*pb.CardListResponse, error) {
			return nil, resperr
		})

		_, err = server.CardList(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("card write", func(t *testing.T) {
		_, err := server.CardWrite(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "card write error")
		server.cardWriteHandler = card.GRPCWriteHandler(func(ctx context.Context, in *pb.CardWriteRequest) (*empty.Empty, error) {
			return nil, resperr
		})

		_, err = server.CardWrite(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("card read", func(t *testing.T) {
		_, err := server.CardRead(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "card read error")
		server.cardReadHandler = card.GRPCReadHandler(func(ctx context.Context, crr *pb.CardReadRequest) (*pb.CardReadResponse, error) {
			return nil, resperr
		})

		_, err = server.CardRead(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("card update", func(t *testing.T) {
		_, err := server.CardUpdate(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "card update error")
		server.cardUpdateHandler = card.GRPCUpdateHandler(func(ctx context.Context, cur *pb.CardUpdateRequest) (*empty.Empty, error) {
			return nil, resperr
		})

		_, err = server.CardUpdate(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("card delete", func(t *testing.T) {
		_, err := server.CardDelete(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "card delete error")
		server.cardDeleteHandler = card.GRPCDeleteHandler(func(ctx context.Context, in *pb.CardDelRequest) (*empty.Empty, error) {
			return nil, resperr
		})

		_, err = server.CardDelete(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	// note

	t.Run("note list", func(t *testing.T) {
		_, err := server.NoteList(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "note list error")
		server.noteListHandler = note.GRPCListHandler(func(ctx context.Context, in *empty.Empty) (*pb.NoteListResponse, error) {
			return nil, resperr
		})

		_, err = server.NoteList(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("note write", func(t *testing.T) {
		_, err := server.NoteWrite(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "note write error")
		server.noteWriteHandler = note.GRPCWriteHandler(func(ctx context.Context, in *pb.NoteWriteRequest) (*empty.Empty, error) {
			return nil, resperr
		})

		_, err = server.NoteWrite(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("note read", func(t *testing.T) {
		_, err := server.NoteRead(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "note read error")
		server.noteReadHandler = note.GRPCReadHandler(func(ctx context.Context, nrr *pb.NoteReadRequest) (*pb.NoteReadResponse, error) {
			return nil, resperr
		})

		_, err = server.NoteRead(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("note update", func(t *testing.T) {
		_, err := server.NoteUpdate(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "note update error")
		server.noteUpdateHandler = note.GRPCUpdateHandler(func(ctx context.Context, nur *pb.NoteUpdateRequest) (*empty.Empty, error) {
			return nil, resperr
		})

		_, err = server.NoteUpdate(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("note delete", func(t *testing.T) {
		_, err := server.NoteDelete(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "note delete error")
		server.noteDeleteHandler = note.GRPCDeleteHandler(func(ctx context.Context, in *pb.NoteDelRequest) (*empty.Empty, error) {
			return nil, resperr
		})

		_, err = server.NoteDelete(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	// binary

	t.Run("binary list", func(t *testing.T) {
		_, err := server.BinaryList(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "binary list error")
		server.binaryListHandler = binary.GRPCListHandler(func(ctx context.Context, in *empty.Empty) (*pb.BinaryListResponse, error) {
			return nil, resperr
		})

		_, err = server.BinaryList(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("binary write", func(t *testing.T) {
		_, err := server.BinaryWrite(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "binary write error")
		server.binaryWriteHandler = binary.GRPCWriteHandler(func(ctx context.Context, in *pb.BinaryWriteRequest) (*pb.BinaryWriteResponse, error) {
			return nil, resperr
		})

		_, err = server.BinaryWrite(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("binary read", func(t *testing.T) {
		_, err := server.BinaryRead(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "binary read error")
		server.binaryReadHandler = binary.GRPCReadHandler(func(ctx context.Context, brr *pb.BinaryReadRequest) (*pb.BinaryReadResponse, error) {
			return nil, resperr
		})

		_, err = server.BinaryRead(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("binary update", func(t *testing.T) {
		_, err := server.BinaryUpdate(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "binary update error")
		server.binaryUpdateHandler = binary.GRPCUpdateHandler(func(ctx context.Context, bur *pb.BinaryUpdateRequest) (*empty.Empty, error) {
			return nil, resperr
		})

		_, err = server.BinaryUpdate(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("binary delete", func(t *testing.T) {
		_, err := server.BinaryDelete(ctx, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "binary delete error")
		server.binaryDeleteHandler = binary.GRPCDeleteHandler(func(ctx context.Context, in *pb.BinaryDelRequest) (*empty.Empty, error) {
			return nil, resperr
		})

		_, err = server.BinaryDelete(ctx, nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("binary upload", func(t *testing.T) {
		err := server.BinaryUpload(nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "binary upload error")
		server.binaryUploadHandler = binary.GRPCUploadHandler(func(us pb.GophKeeper_BinaryUploadServer) error {
			return resperr
		})

		err = server.BinaryUpload(nil)
		require.ErrorIs(t, err, resperr)
	})

	t.Run("binary download", func(t *testing.T) {
		err := server.BinaryDownload(nil, nil)
		require.Error(t, err)

		resperr := status.Error(codes.Internal, "binary download error")
		server.binaryDownloadHandler = binary.GRPCDownloadHandler(func(req *pb.BidaryDownloadRequest, ds pb.GophKeeper_BinaryDownloadServer) error {
			return resperr
		})

		err = server.BinaryDownload(nil, nil)
		require.ErrorIs(t, err, resperr)
	})

}
