// Package client - grpc клиент
package client

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"github.com/golang/protobuf/ptypes/empty"
)

type Client struct {
	conn       *grpc.ClientConn
	client     pb.GophKeeperClient
	userTokens map[string]string
	userName   string
}

func NewClient(addr string) (*Client, error) {
	var (
		client Client
		err    error
	)

	// устанавливаем соединение с сервером
	client.conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(echoInterceptor))
	if err != nil {
		return nil, err
	}

	client.userTokens = make(map[string]string, 1)

	// Получаем переменную интерфейсного типа UserClient,
	// через которую будем отправлять сообщения
	client.client = pb.NewGophKeeperClient(client.conn)
	return &client, nil
}

func echoInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	// выполняем действия перед вызовом метода
	start := time.Now()

	logger.Debug("request",
		"method", method,
		"data", req)

	// вызов RPC-метод
	err := invoker(ctx, method, req, reply, cc, opts...)

	// выводим действия после вызова метода
	if err != nil {
		logger.Debug(err.Error())
	} else {
		logger.Debug("response",
			"duration", time.Since(start),
			"data", reply)
	}
	return err
}

func (c *Client) withUserToken(ctx context.Context, userNane string) context.Context {
	md := metadata.New(map[string]string{
		"token": c.userTokens[userNane],
	})
	return metadata.NewOutgoingContext(ctx, md)
}

func (c *Client) withToken(ctx context.Context) context.Context {
	return c.withUserToken(ctx, c.userName)
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetUser() string {
	return c.userName
}

func (c *Client) SetUser(name string) error {
	if name == "" {
		return fmt.Errorf("username is empty")
	}
	_, ok := c.userTokens[name]
	if ok {
		c.userName = name
		return nil
	}
	return fmt.Errorf("user '%s' not found", name)
}

func (c *Client) GetUsers() []string {
	res := make([]string, 0, len(c.userTokens))
	for u := range c.userTokens {
		res = append(res, u)
	}
	return res
}

func (c *Client) Ping() error {
	_, err := c.client.Ping(context.Background(), &emptypb.Empty{})
	return err
}

func (c *Client) Login(login, passwd string) error {
	req := pb.LoginRequest{
		Login:    login,
		Password: passwd,
	}
	resp, err := c.client.Login(context.Background(), &req)
	if err == nil {
		c.userName = login
		c.userTokens[login] = resp.Token
	}
	return err
}

func (c *Client) Registration(login, passwd string) error {
	req := pb.RegisterRequest{
		Login:    login,
		Password: passwd,
	}
	resp, err := c.client.Register(context.Background(), &req)
	if err == nil {
		c.userName = login
		c.userTokens[login] = resp.Token
	}
	return err
}

func (c *Client) List() (*pb.ListResponse, error) {
	ctx := c.withToken(context.Background())
	return c.client.List(ctx, &empty.Empty{})
}

// Cards //

func (c *Client) CardList() ([]string, error) {
	ctx := c.withToken(context.Background())
	resp, err := c.client.CardList(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.Names, nil
}

func (c *Client) CardWrite(in *pb.CardWriteRequest) error {
	ctx := c.withToken(context.Background())
	_, err := c.client.CardWrite(ctx, in)
	return err
}

func (c *Client) CardRead(in *pb.CardReadRequest) (*pb.CardReadResponse, error) {
	ctx := c.withToken(context.Background())
	return c.client.CardRead(ctx, in)
}

func (c *Client) CardUpdate(name string, in *pb.CardWriteRequest) error {
	ctx := c.withToken(context.Background())
	pass, err := c.client.CardRead(ctx, &pb.CardReadRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	req := pb.CardUpdateRequest{
		Id:    pass.Id,
		Write: in,
	}
	_, err = c.client.CardUpdate(ctx, &req)
	return err
}

func (c *Client) CardDelete(in *pb.CardDelRequest) error {
	ctx := c.withToken(context.Background())
	_, err := c.client.CardDelete(ctx, in)
	return err
}

// Notes //

func (c *Client) NoteList() ([]string, error) {
	ctx := c.withToken(context.Background())
	resp, err := c.client.NoteList(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.Names, nil
}

func (c *Client) NoteWrite(in *pb.NoteWriteRequest) error {
	ctx := c.withToken(context.Background())
	_, err := c.client.NoteWrite(ctx, in)
	return err
}

func (c *Client) NoteRead(in *pb.NoteReadRequest) (*pb.NoteReadResponse, error) {
	ctx := c.withToken(context.Background())
	return c.client.NoteRead(ctx, in)
}

func (c *Client) NoteUpdate(name string, in *pb.NoteWriteRequest) error {
	ctx := c.withToken(context.Background())
	pass, err := c.client.NoteRead(ctx, &pb.NoteReadRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	req := pb.NoteUpdateRequest{
		Id:    pass.Id,
		Write: in,
	}
	_, err = c.client.NoteUpdate(ctx, &req)
	return err
}

func (c *Client) NoteDelete(in *pb.NoteDelRequest) error {
	ctx := c.withToken(context.Background())
	_, err := c.client.NoteDelete(ctx, in)
	return err
}

// Binaries //

func (c *Client) BinaryList() ([]string, error) {
	ctx := c.withToken(context.Background())
	resp, err := c.client.BinaryList(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.Names, nil
}

func (c *Client) BinaryWrite(in *pb.BinaryWriteRequest) (int64, error) {
	ctx := c.withToken(context.Background())
	resp, err := c.client.BinaryWrite(ctx, in)
	if err != nil {
		return 0, err
	}
	return resp.Id, nil
}

func (c *Client) BinaryRead(in *pb.BinaryReadRequest) (*pb.BinaryReadResponse, error) {
	ctx := c.withToken(context.Background())
	return c.client.BinaryRead(ctx, in)
}

func (c *Client) BinaryUpdate(id int64, binID int64, in *pb.BinaryWriteRequest) error {
	ctx := c.withToken(context.Background())
	req := pb.BinaryUpdateRequest{
		Id:    id,
		BinId: binID,
		Write: in,
	}
	_, err := c.client.BinaryUpdate(ctx, &req)
	return err
}

func (c *Client) BinaryDelete(in *pb.BinaryDelRequest) error {
	ctx := c.withToken(context.Background())
	_, err := c.client.BinaryDelete(ctx, in)
	return err
}

func (c *Client) BinaryUpload(id int64, r io.Reader) error {
	ctx := c.withToken(context.Background())
	client, err := c.client.BinaryUpload(ctx)
	if err != nil {
		return err
	}

	chSize := 4096
	upload := pb.BinaryUplodStream{
		Chunk: make([]byte, chSize),
	}

	for err == nil {
		var n int
		n, err = r.Read(upload.Chunk)
		if n > 0 && (err == nil || err == io.EOF) {
			upload.Id = id
			upload.Chunk = upload.Chunk[:n]
			err = client.Send(&upload)
		}
	}
	if err == io.EOF {
		err = nil
	}
	if _, e := client.CloseAndRecv(); e != nil && err == nil {
		err = e
	}
	return err
}

func (c *Client) BinaryDownload(id int64, w io.Writer) error {

	ctx := c.withToken(context.Background())
	req := pb.BidaryDownloadRequest{
		Id: id,
	}

	download, err := c.client.BinaryDownload(ctx, &req)
	if err != nil {
		return err
	}

	for err == nil {
		var stream *pb.BinaryDownloadStream
		stream, err = download.Recv()
		if err == nil {
			logger.Debug("download", "err", err, "stream", stream)
			_, err = w.Write(stream.Chunk)
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}

// Password //

func (c *Client) PasswordList() ([]string, error) {
	ctx := c.withToken(context.Background())
	resp, err := c.client.PasswordList(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.Names, nil
}

func (c *Client) PasswordWrite(in *pb.PasswordWriteRequest) error {
	ctx := c.withToken(context.Background())
	_, err := c.client.PasswordWrite(ctx, in)
	return err
}

func (c *Client) PasswordRead(in *pb.PasswordReadRequest) (*pb.PasswordReadResponse, error) {
	ctx := c.withToken(context.Background())
	return c.client.PasswordRead(ctx, in)
}

func (c *Client) PasswordUpdate(name string, in *pb.PasswordWriteRequest) error {
	ctx := c.withToken(context.Background())
	pass, err := c.client.PasswordRead(ctx, &pb.PasswordReadRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	req := pb.PasswordUpdateRequest{
		Id:    pass.Id,
		Write: in,
	}
	_, err = c.client.PasswordUpdate(ctx, &req)
	return err
}

func (c *Client) PasswordDelete(in *pb.PasswordDelRequest) error {
	ctx := c.withToken(context.Background())
	_, err := c.client.PasswordDelete(ctx, in)
	return err
}
