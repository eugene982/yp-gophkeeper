package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/golang/protobuf/ptypes/empty"
)

const (
	addr           = ":8080"
	requestTimeout = time.Second * 10
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(clientInterceptor))
	if err != nil {
		return err
	}
	defer conn.Close()

	// Получаем переменную интерфейсного типа UserClient,
	// через которую будем отправлять сообщения
	c := pb.NewGophKeeperClient(conn)

	_, err = Ping(c)
	if err != nil {
		return fmt.Errorf("ping error: %w", err)
	}

	// ***
	err = Register(c, "", "")
	if err != nil {
		return fmt.Errorf("register error: %w", err)
	}

	return nil
}

func clientInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	// выполняем действия перед вызовом метода
	start := time.Now()

	// вызов RPC-метод
	err := invoker(ctx, method, req, reply, cc, opts...)

	// выводим действия после вызова метода
	if err != nil {
		log.Printf("[ERROR]: %s, %v", method, err)
	} else {
		log.Printf("[INFO]: %v, %s, %v", time.Since(start), method, reply)
	}
	return err
}

func Ping(c pb.GophKeeperClient) (string, error) {
	resp, err := c.Ping(context.Background(), &empty.Empty{})
	if err != nil {
		return "", err
	}
	return resp.Message, nil
}

func Register(c pb.GophKeeperClient, login, passwd string) error {
	req := pb.RegisterRequest{
		Login:    login,
		Password: passwd,
	}
	_, err := c.Register(context.Background(), &req)
	return err
}
