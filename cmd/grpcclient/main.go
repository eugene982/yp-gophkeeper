package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/golang/protobuf/ptypes/empty"
)

var (
	debugMode bool
	userName  string // имя пользователя
	userToken string // токен пользователя
	srvAddr   string // адрес сервера
	client    pb.GophKeeperClient

	// обработчики комманд
	handlers = map[string](func() error){
		"help":  helpCmd,
		"ping":  pingCmd,
		"reg":   regCmd,
		"ls":    lsCmd,
		"list":  lsCmd,
		"login": loginCmd,
	}
)

func main() {

	flag.StringVar(&srvAddr, "a", ":28000", "gophkeeper server address")
	flag.BoolVar(&debugMode, "debug", true, "debug mode")

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(srvAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(echoInterceptor))
	if err != nil {
		return err
	}
	defer conn.Close()

	// Получаем переменную интерфейсного типа UserClient,
	// через которую будем отправлять сообщения
	client = pb.NewGophKeeperClient(conn)

	var (
		cmd string = "help"
	)

	for mainLop(cmd) {
		prompt()

		_, err := fmt.Scanln(&cmd)
		if err != nil {
			log.Println("fmt.Scanln err:", err)
			if err == io.EOF {
				return nil
			}
			fmt.Fprintln(os.Stderr, "wrong command")
			cmd = "help"
			continue
		}

	}
	return nil
}

func mainLop(cmd string) bool {
	if cmd == "q" || cmd == "quit" || cmd == "exit" {
		return false
	} else if cmd == "" {
		return true
	}

	fn, ok := handlers[cmd]
	if !ok {
		fmt.Fprintln(os.Stderr, "wrong command")
		helpCmd()
		return true
	}

	err := fn()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return true
}

func prompt() {
	if userName == "" {
		fmt.Printf("%s>", srvAddr)
	} else {
		fmt.Printf("%s@%s>", userName, srvAddr)
	}
}

func helpCmd() error {
	fmt.Println(`	"help"       - вывод справки по командам 
	"quit" ("q") - выход из программы"
	"ping"       - проверка соединение (пинг)
	"reg"        - регистрация (создание) нового пользователя
	"login"      - вход`)
	return nil
}

func pingCmd() error {
	resp, err := client.Ping(context.Background(), &empty.Empty{})
	if err != nil {
		return err
	}
	fmt.Println(resp.Message)
	return nil
}

func regCmd() error {
	var login, passwd string

	fmt.Print("\tlogin:")
	_, err := fmt.Scanln(&login)
	if err != nil {
		return err
	}

	fmt.Print("\tpassword:")
	_, err = fmt.Scanln(&passwd)
	if err != nil {
		return err
	}

	req := pb.RegisterRequest{
		Login:    login,
		Password: passwd,
	}
	resp, err := client.Register(context.Background(), &req)
	if err != nil {
		return err
	}
	userToken = resp.Token
	userName = login
	fmt.Println("\tOK")
	return nil
}

func loginCmd() error {
	var login, passwd string

	fmt.Print("\tlogin:")
	_, err := fmt.Scanln(&login)
	if err != nil {
		return err
	}

	fmt.Print("\tpassword:")
	_, err = fmt.Scanln(&passwd)
	if err != nil {
		return err
	}

	req := pb.LoginRequest{
		Login:    login,
		Password: passwd,
	}
	resp, err := client.Login(context.Background(), &req)
	if err != nil {
		return err
	}
	userToken = resp.Token
	userName = login
	fmt.Println("\tOK")
	return nil
}

func lsCmd() error {

	ctx := withToken(context.Background())
	resp, err := client.List(ctx, &empty.Empty{})
	if err != nil {
		return nil
	}

	fmt.Println("\tNotes    :", resp.NotesCount)
	fmt.Println("\tCards    :", resp.CardsCount)
	fmt.Println("\tPasswords:", resp.PasswordsCount)

	return nil
}

func withToken(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{
		"token": userToken,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

func echoInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	if !debugMode {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	// выполняем действия перед вызовом метода
	start := time.Now()

	log.Printf("[REQ]: %v, %s, %v", start, method, req)

	// вызов RPC-метод
	err := invoker(ctx, method, req, reply, cc, opts...)

	// выводим действия после вызова метода
	if err != nil {
		log.Printf("[ERROR]: %s", err)
	} else {
		log.Printf("[RESP]: %v, %v", time.Since(start), reply)
	}
	return err
}
