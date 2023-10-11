package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
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
		return err
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
		log.Printf("[ERROR] %s,%v", method, err)
	} else {
		log.Printf("[INFO], %s,%v", method, time.Since(start))
	}
	return err
}

func Ping(c pb.GophKeeperClient) (string, error) {
	ctx, candel := context.WithTimeout(context.Background(), requestTimeout)
	defer candel()
	resp, err := c.Ping(ctx, nil)
	return resp.Message, err
}
