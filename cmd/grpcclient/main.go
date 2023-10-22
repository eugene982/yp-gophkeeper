package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
	currTable string // текущая таблица

	srvAddr string // адрес сервера
	client  pb.GophKeeperClient

	errUnauthenticated = errors.New("unauthenticated")

	conReader = consoleReader{bufio.NewReader(os.Stdin)}

	// обработчики комманд
	handlers = map[string](func([]string) error){
		"help":  helpCmd,
		"ping":  pingCmd,
		"reg":   regCmd,
		"ls":    lsCmd,
		"list":  lsCmd,
		"login": loginCmd,

		"new": newCmd,
		"get": getCmd,
		"del": delCmd,
		"upd": updCmd,

		"passwords":   passwordsCmd,
		"lspasswords": lsPasswordCmd,
		"newpassword": newPasswordCmd,

		"cards": lsCardCmd,
		"lscrd": lsCardCmd,

		"notes": lsNoteCmd,
		"lsnt":  lsNoteCmd,
	}
)

type consoleReader struct {
	reader *bufio.Reader
}

func (c consoleReader) ReadLine() (string, error) {
	s, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

func (c consoleReader) ReadFields() ([]string, error) {
	str, err := c.ReadLine()
	if err != nil {
		return nil, err
	}
	return strings.Split(str, " "), nil
}

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
		cmd     string   = "help"
		cmdArgs []string = nil
	)

	for mainLop(cmd, cmdArgs) {
		prompt()
		cmdArgs = nil

		args, err := conReader.ReadFields()
		if err != nil {
			log.Println("fmt.Scanln err:", err)
			if err == io.EOF {
				return nil
			}
			fmt.Fprintln(os.Stderr, "неверная команда")
			cmd = "help"
			continue
		} else if len(args) == 0 {
			cmd = "help"
			continue
		}
		cmd = args[0]
		cmdArgs = args[1:]
	}
	return nil
}

func mainLop(cmd string, args []string) bool {
	if cmd == "q" || cmd == "quit" || cmd == "exit" {
		return false
	} else if cmd == "" {
		return true
	}

	fn, ok := handlers[cmd]
	if !ok {
		fmt.Fprintln(os.Stderr, "неверная комманда: ", cmd)
		helpCmd(nil)
		return true
	}

	err := fn(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return true
}

func prompt() {
	if userName == "" {
		fmt.Printf("%s>", srvAddr)
	} else if currTable == "" {
		fmt.Printf("%s@%s>", userName, srvAddr)
	} else {
		fmt.Printf("%s@%s/%s>", userName, srvAddr, currTable)
	}
}

func helpCmd(args []string) error {
	switch currTable {
	case "passwords":
		return helpPasswordCmd(args)
	}

	fmt.Println(`	
	"help"              - вывод справки по командам 
	"quit" ("q")        - выход из программы"
	"ping"              - проверка соединение (пинг)
	"reg" [user pass]   - регистрация (создание) нового пользователя
	"login" [user pass] - вход
	"passwords"         - управление хранилищем паролей
	"cards"             - управление хранилищем карт
	"notes"             - управление хранилищем заметок
	`)
	return nil
}

func pingCmd([]string) error {
	resp, err := client.Ping(context.Background(), &empty.Empty{})
	if err != nil {
		return err
	}
	fmt.Println(resp.Message)
	return nil
}

func regCmd(args []string) (err error) {
	var login, passwd string

	if len(args) > 0 {
		login = args[0]
	} else {
		fmt.Print("\tlogin:")
		login, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	if len(args) > 1 {
		passwd = args[1]
	} else {
		fmt.Print("\tpassword:")
		passwd, err = conReader.ReadLine()
		if err != nil {
			return err
		}
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

func loginCmd(args []string) (err error) {
	var login, passwd string

	if len(args) > 0 {
		login = args[0]
	} else {
		fmt.Print("\tlogin:")
		login, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	if len(args) > 1 {
		passwd = args[1]
	} else {
		fmt.Print("\tpassword:")
		passwd, err = conReader.ReadLine()
		if err != nil {
			return err
		}
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

func lsCmd(args []string) error {

	switch currTable {
	case "passwords":
		return lsPasswordCmd(args)
	}

	ctx := withToken(context.Background())
	resp, err := client.List(ctx, &empty.Empty{})
	if err != nil {
		return nil
	}

	fmt.Println("\t    Notes:", resp.NotesCount)
	fmt.Println("\t    Cards:", resp.CardsCount)
	fmt.Println("\tPasswords:", resp.PasswordsCount)

	return nil
}

func newCmd(args []string) error {
	switch currTable {
	case "passwords":
		return newPasswordCmd(args)
	}
	return fmt.Errorf("выберите таблицу: passwords, cards, notes...")
}

func getCmd(args []string) error {
	switch currTable {
	case "passwords":
		return getPasswordCmd(args)
	}
	return fmt.Errorf("выберите таблицу: passwords, cards, notes...")
}

func delCmd(args []string) error {
	switch currTable {
	case "passwords":
		return delPasswordCmd(args)
	}
	return fmt.Errorf("выберите таблицу: passwords, cards, notes...")
}

func updCmd(args []string) error {
	switch currTable {
	case "passwords":
		return updPasswordCmd(args)
	}
	return errors.New("выберите таблицу: passwords, cards, notes...")
}

func lsCardCmd([]string) error {

	ctx := withToken(context.Background())
	resp, err := client.CardList(ctx, &empty.Empty{})
	if err != nil {
		return nil
	}

	fmt.Println("\tCards:", strings.Join(resp.Names, "; "))
	fmt.Println("\tCount:", len(resp.Names))

	return nil
}

func lsNoteCmd([]string) error {

	ctx := withToken(context.Background())
	resp, err := client.NoteList(ctx, &empty.Empty{})
	if err != nil {
		return nil
	}

	fmt.Println("\tNotes:", strings.Join(resp.Names, "; "))
	fmt.Println("\tCount:", len(resp.Names))

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
