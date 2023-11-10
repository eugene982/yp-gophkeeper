package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/c-bata/go-prompt"

	"github.com/eugene982/yp-gophkeeper/cmd/grpcclient/client"
	"github.com/eugene982/yp-gophkeeper/cmd/grpcclient/command"
	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
)

var (
	gkeeperClient *client.Client
)

var (
	logLevel                             string
	serverAddress                        string
	buildVersion, buildDate, buildCommit string
)

func main() {

	fmt.Println("gophkeeper CLI")
	if buildVersion != "" {
		fmt.Println("version:", buildVersion, buildDate, buildCommit)
	}
	fmt.Println(`"exit" или "Ctrl-D" чтобы выйти из программы`)

	flag.StringVar(&logLevel, "l", "error", "log level")
	flag.StringVar(&serverAddress, "a", ":28000", "gophkeeper server addres")
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	err := logger.Initialize(logLevel)
	if err != nil {
		return err
	}

	gkeeperClient, err = client.NewClient(serverAddress)
	if err != nil {
		return err
	}

	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("interactive shell client"),
		prompt.OptionLivePrefix(livePrefix),
		prompt.OptionInputTextColor(prompt.Yellow),
		//prompt.OptionSetExitCheckerOnInput(func(in string, breakline bool) bool { return in == "exit" }),
	)

	p.Run()
	return nil
}

func executor(line string) {
	var cmd *command.Command

	commands := strings.Split(strings.TrimSpace(line), " ")
	args := commands[1:]

	switch commands[0] {
	case "exit", "quit":
		cmd = newExitCmd(args)
	case "info":
		cmd = newInfoCmd(args)
	case "ping":
		cmd = newPingCmd(args)
	case "login":
		cmd = newLoginCmd(args)
	case "reg":
		cmd = newRegCmd(args)
	case "user":
		cmd = newUserCmd(args)
	case "ls", "list":
		cmd = newListCmd(args)
	case "card":
		cmd = newCardsCmd(args)
	case "note":
		cmd = newNotesCmd(args)
	case "file":
		cmd = newFilesCmd(args)
	case "password":
		cmd = newPasswordsCmd(args)
	default:
		fmt.Println("неизвестная команда:", line)
		return
	}
	cmd.Exec()
}

func completer(d prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest

	words := strings.Split(strings.TrimLeft(d.Text, " "), " ")

	//logger.Debug("completer", "words", words, "count", wordsCount)

	switch len(words) {
	case 0, 1:
		s = []prompt.Suggest{
			{Text: "exit", Description: "выход из клиента"},
			{Text: "info", Description: "вывод информации о сборке"},
			{Text: "ping", Description: "проверка соединения"},
			{Text: "login", Description: "[user password] авторизация пользователя"},
			{Text: "reg", Description: "[user password] регистрация нового пользователя"},
			{Text: "user", Description: "[name] выбор авторизированного польтзователя"},

			{Text: "list", Description: "список хранимых данных"},
			{Text: "ls", Description: "список хранимых данных"},

			{Text: "password", Description: "работа с хранилищем паролей"},
			{Text: "note", Description: "работа с хранилищем заметок"},
			{Text: "card", Description: "работа с хранилищем карт"},
			{Text: "file", Description: "работа с хранилищем файлов"},
		}
	case 2:
		switch words[0] {
		case "user":
			for _, u := range gkeeperClient.GetUsers() {
				s = append(s, prompt.Suggest{Text: u})
			}
		case "password", "note", "card", "file":
			s = []prompt.Suggest{
				{Text: "ls", Description: "показать список"},
				{Text: "get", Description: "прочитать данные из хранилища"},
				{Text: "new", Description: "добавить данные в хранилище"},
				{Text: "upd", Description: "обновить данные"},
				{Text: "del", Description: "удалить из хранилища"},
			}
		}
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func livePrefix() (prefix string, useLivePrefix bool) {
	useLivePrefix = true
	if gkeeperClient.GetUser() != "" {
		prefix = gkeeperClient.GetUser() + "@" + serverAddress + "> "
	} else {
		prefix = serverAddress + "> "
	}
	return
}

func newPingCmd(args []string) *command.Command {
	return command.New(func(m map[string]string) error {
		logger.Debug("ping", "args", args)
		err := gkeeperClient.Ping()
		if err == nil {
			fmt.Println("OK")
		}
		return err
	},
		args)
}

func newExitCmd(args []string) *command.Command {
	return command.New(func(m map[string]string) error {
		if err := gkeeperClient.Close(); err != nil {
			logger.Errorf("close error: %w", err)
		}
		fmt.Println("Bye!")
		os.Exit(0)
		return nil
	},
		args)
}

func newInfoCmd(args []string) *command.Command {
	return command.New(func(m map[string]string) error {
		fmt.Println("version:", buildVersion)
		fmt.Println("date:   ", buildDate)
		fmt.Println("commit: ", buildCommit)
		return nil
	},
		args)
}

func newRegCmd(args []string) *command.Command {
	return command.New(func(m map[string]string) error {
		login := m["login"]
		passwd := m["password"]
		return gkeeperClient.Registration(login, passwd)
	},
		args,
		"login", "password")
}

func newLoginCmd(args []string) *command.Command {
	return command.New(func(m map[string]string) error {
		login := m["login"]
		passwd := m["password"]
		return gkeeperClient.Login(login, passwd)
	},
		args,
		"login", "password")
}

func newUserCmd(args []string) *command.Command {
	return command.New(func(m map[string]string) error {
		return gkeeperClient.SetUser(m["name"])
	},
		args,
		"name")
}

func newListCmd(args []string) *command.Command {
	return command.New(func(m map[string]string) error {
		resp, err := gkeeperClient.List()
		if err == nil {
			fmt.Println("cards:", resp.CardsCount)
			fmt.Println("files:", resp.BinariesCount)
			fmt.Println("notes:", resp.NotesCount)
			fmt.Println("passwords:", resp.PasswordsCount)
		}
		return err
	}, args)
}

// newCardCmd - обработчики команда работы с картами
func newCardsCmd(args []string) *command.Command {
	var (
		subcmd  string
		subargs []string
	)
	errCmd := command.New(func(map[string]string) error {
		return fmt.Errorf("неизвестная команда: %s", strings.Join(args, " "))
	}, nil)

	if len(args) > 0 {
		subcmd = args[0]
		subargs = args[1:]
	}

	switch subcmd {
	case "", "ls", "list":
		return command.New(func(m map[string]string) error {
			names, err := gkeeperClient.CardList()
			if err != nil {
				return err
			} else if len(names) == 0 {
				fmt.Println("нет сохраненных карт")
			} else {
				fmt.Println(strings.Join(names, "\n"))
			}
			return nil
		}, subargs)
	case "new":
		return command.New(func(fields map[string]string) error {
			in := pb.CardWriteRequest{
				Name:   fields["name"],
				Number: fields["number"],
				Pin:    fields["pin"],
				Notes:  fields["notes"],
			}
			return gkeeperClient.CardWrite(&in)
		}, subargs, "name", "number", "pin", "notes")
	case "get":
		return command.New(func(fields map[string]string) error {
			in := pb.CardReadRequest{
				Name: fields["name"],
			}
			resp, err := gkeeperClient.CardRead(&in)
			if err == nil {
				fmt.Println("name:", resp.Name)
				fmt.Println("number:", resp.Number)
				fmt.Println("pin:", resp.Pin)
				fmt.Println("notes:", resp.Notes)
			}
			return err
		}, subargs, "name")
	case "upd":
		return command.New(func(fields map[string]string) error {
			in := pb.CardWriteRequest{
				Name:   fields["new name"],
				Number: fields["new number"],
				Pin:    fields["new pin"],
				Notes:  fields["new notes"],
			}
			return gkeeperClient.CardUpdate(fields["name"], &in)
		}, subargs, "name", "new name", "new number", "new pin", "new notes")
	case "del":
		return command.New(func(fields map[string]string) error {
			in := pb.CardDelRequest{
				Name: fields["name"],
			}
			return gkeeperClient.CardDelete(&in)
		}, subargs, "name")
	}
	return errCmd
}

// newNotesCmd - обработчики команда работы с заметками
func newNotesCmd(args []string) *command.Command {
	var (
		subcmd  string
		subargs []string
	)
	if len(args) > 0 {
		subcmd = args[0]
		subargs = args[1:]
	}

	switch subcmd {
	case "", "ls", "list":
		return command.New(func(m map[string]string) error {
			names, err := gkeeperClient.NoteList()
			if err != nil {
				return err
			} else if len(names) == 0 {
				fmt.Println("нет сохраненных заметок")
			} else {
				fmt.Println(strings.Join(names, "\n"))
			}
			return nil
		}, subargs)

	case "new":
		return command.New(func(fields map[string]string) error {
			in := pb.NoteWriteRequest{
				Name:  fields["name"],
				Notes: fields["notes"],
			}
			return gkeeperClient.NoteWrite(&in)
		}, subargs, "name", "notes")

	case "get":
		return command.New(func(fields map[string]string) error {
			in := pb.NoteReadRequest{
				Name: fields["name"],
			}
			resp, err := gkeeperClient.NoteRead(&in)
			if err == nil {
				fmt.Println("name:", resp.Name)
				fmt.Println("notes:", resp.Notes)
			}
			return err
		}, subargs, "name")

	case "upd":
		return command.New(func(fields map[string]string) error {
			in := pb.NoteWriteRequest{
				Name:  fields["new name"],
				Notes: fields["new notes"],
			}
			return gkeeperClient.NoteUpdate(fields["name"], &in)
		}, subargs, "name", "new name", "new notes")

	case "del":
		return command.New(func(fields map[string]string) error {
			in := pb.NoteDelRequest{
				Name: fields["name"],
			}
			return gkeeperClient.NoteDelete(&in)
		}, subargs, "name")

	default:
		return command.New(func(map[string]string) error {
			return fmt.Errorf("неизвестная команда: %s", strings.Join(args, " "))
		}, nil)
	}
}

// newNotesCmd - обработчики команда работы с заметками
func newFilesCmd(args []string) *command.Command {
	var (
		subcmd  string
		subargs []string
	)
	if len(args) > 0 {
		subcmd = args[0]
		subargs = args[1:]
	}

	switch subcmd {
	case "", "ls", "list":
		return command.New(func(m map[string]string) error {
			names, err := gkeeperClient.BinaryList()
			if err != nil {
				return err
			} else if len(names) == 0 {
				fmt.Println("нет сохраненных файлов")
			} else {
				fmt.Println(strings.Join(names, "\n"))
			}
			return nil
		}, subargs)

	case "new":
		return command.New(func(fields map[string]string) error {
			filename := fields["file"]
			file, err := os.Open(filename)
			if err != nil {
				return err
			}
			fstat, err := file.Stat()
			if err != nil {
				return err
			}
			defer file.Close()

			name := fields["name"]
			if name == "" {
				name = path.Base(filename)
			}

			//резервируем идентификатор под файл
			in := pb.BinaryWriteRequest{
				Name:  name,
				Notes: fields["notes"],
				Size:  fstat.Size(),
			}

			id, err := gkeeperClient.BinaryWrite(&in)
			if err != nil {
				return err
			}
			return gkeeperClient.BinaryUpload(id, file)

		}, subargs, "file", "name", "notes")

	case "get":
		return command.New(func(fields map[string]string) error {
			req := pb.BinaryReadRequest{
				Name: fields["name"],
			}

			resp, err := gkeeperClient.BinaryRead(&req)
			if err != nil {
				return err
			}
			fmt.Println("name:", resp.Name)
			fmt.Println("size:", resp.Size)
			fmt.Println("notes:", resp.Notes)

			filename := fields["save to file"]
			if filename == "" {
				filename = resp.Name
			}
			file, err := os.Create(filename)
			if err != nil {
				return err
			}
			defer file.Close()

			err = gkeeperClient.BinaryDownload(resp.BinId, file)
			return err
		}, subargs, "name", "save to file")

	case "upd":
		return command.New(func(fields map[string]string) error {
			req := pb.BinaryReadRequest{
				Name: fields["name"],
			}

			resp, err := gkeeperClient.BinaryRead(&req)
			if err != nil {
				return err
			}

			upd := pb.BinaryWriteRequest{
				Name:  fields["new name"],
				Notes: fields["new notes"],
				Size:  resp.Size,
			}
			if upd.Name == "" {
				upd.Name = resp.Name
			}
			if upd.Notes == "" {
				upd.Notes = resp.Notes
			}

			var (
				file  *os.File
				binID int64
			)
			filename := fields["new file"]

			if filename != "" {
				binID = resp.BinId
				file, err = os.Open(filename)
				if err != nil {
					return err
				}
				fstat, err := file.Stat()
				if err != nil {
					return err
				}
				defer file.Close()
				upd.Size = fstat.Size()
			}

			err = gkeeperClient.BinaryUpdate(resp.Id, binID, &upd)
			if err != nil {
				return err
			}

			if file != nil {
				return gkeeperClient.BinaryUpload(resp.BinId, file)
			}
			return nil

		}, subargs, "name", "new file", "new name", "new notes")

	case "del":
		return command.New(func(fields map[string]string) error {
			in := pb.BinaryDelRequest{
				Name: fields["name"],
			}
			return gkeeperClient.BinaryDelete(&in)
		}, subargs, "name")

	default:
		return command.New(func(map[string]string) error {
			return fmt.Errorf("неизвестная команда: %s", strings.Join(args, " "))
		}, nil)
	}
}

// newCardCmd - обработчики команда работы с паролями
func newPasswordsCmd(args []string) *command.Command {
	var (
		subcmd  string
		subargs []string
	)
	errCmd := command.New(func(map[string]string) error {
		return fmt.Errorf("неизвестная команда: %s", strings.Join(args, " "))
	}, nil)

	if len(args) > 0 {
		subcmd = args[0]
		subargs = args[1:]
	}

	switch subcmd {
	// Список паролей
	case "", "ls", "list":
		return command.New(func(m map[string]string) error {
			names, err := gkeeperClient.PasswordList()
			if err != nil {
				return err
			} else if len(names) == 0 {
				fmt.Println("нет сохраненных паролей")
			} else {
				fmt.Println(strings.Join(names, "\n"))
			}
			return nil
		}, subargs)
	// Создание нового
	case "new":
		return command.New(func(fields map[string]string) error {
			in := pb.PasswordWriteRequest{
				Name:     fields["name"],
				Username: fields["username"],
				Password: fields["password"],
				Notes:    fields["notes"],
			}
			return gkeeperClient.PasswordWrite(&in)
		}, subargs, "name", "username", "password", "notes")
	// получение
	case "get":
		return command.New(func(fields map[string]string) error {
			in := pb.PasswordReadRequest{
				Name: fields["name"],
			}
			resp, err := gkeeperClient.PasswordRead(&in)
			if err == nil {
				fmt.Println("name:", resp.Name)
				fmt.Println("username:", resp.Username)
				fmt.Println("password:", resp.Password)
				fmt.Println("notes:", resp.Notes)
			}
			return err
		}, subargs, "name")
	// обновление
	case "upd":
		return command.New(func(fields map[string]string) error {
			in := pb.PasswordWriteRequest{
				Name:     fields["new name"],
				Username: fields["new username"],
				Password: fields["new password"],
				Notes:    fields["new notes"],
			}
			return gkeeperClient.PasswordUpdate(fields["name"], &in)
		}, subargs, "name", "new name", "new username", "new password", "new notes")
	// Удаление
	case "del":
		return command.New(func(fields map[string]string) error {
			in := pb.PasswordDelRequest{
				Name: fields["name"],
			}
			return gkeeperClient.PasswordDelete(&in)
		}, subargs, "name")
	}
	return errCmd
}
