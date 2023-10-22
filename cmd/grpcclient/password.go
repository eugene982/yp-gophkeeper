package main

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/golang/protobuf/ptypes/empty"
)

func passwordsCmd([]string) error {
	if userName == "" {
		return errUnauthenticated
	}
	currTable = "passwords"
	return nil
}

func helpPasswordCmd(args []string) error {
	fmt.Println(`	
	"ls"         - список сохраненных паролей
	"new [name]" - создание нового пароля
	"get [name]" - получение из хранилища
	"upd [name]" - обновление данных в хранилище
	"del [name]" - получение пароля из хранилища`)
	return nil
}

func lsPasswordCmd([]string) error {

	ctx := withToken(context.Background())
	resp, err := client.PasswordList(ctx, &empty.Empty{})
	if err != nil {
		return nil
	}

	fmt.Println("\tPasswords:", strings.Join(resp.Names, "; "))
	fmt.Println("\tCount: ", len(resp.Names))

	return nil
}

func newPasswordCmd(args []string) (err error) {
	var req pb.PasswordWriteRequest

	if len(args) > 0 {
		req.Name = args[0]
	} else {
		fmt.Print("\t    Name: ")
		req.Name, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	fmt.Print("\tUsername: ")
	req.Username, err = conReader.ReadLine()
	if err != nil {
		return err
	}

	fmt.Print("\tPassword: ")
	req.Password, err = conReader.ReadLine()
	if err != nil {
		return err
	}

	fmt.Print("\t   Notes: ")
	req.Notes, err = conReader.ReadLine()
	if err != nil {
		return err
	}

	ctx := withToken(context.Background())
	_, err = client.PasswordWrite(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Println("\tДобавлен.")
	return nil
}

func getPasswordCmd(args []string) (err error) {
	var req pb.PasswordReadRequest

	if len(args) > 0 {
		req.Name = args[0]
	} else {
		fmt.Print("\t    Name: ")
		req.Name, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	ctx := withToken(context.Background())

	resp, err := client.PasswordRead(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Println("\tUsername:", resp.Username)
	fmt.Println("\tPassword:", resp.Password)
	fmt.Println("\t   Notes:", resp.Notes)
	return nil
}

func delPasswordCmd(args []string) (err error) {
	var req pb.PasswordDelRequest

	if len(args) > 0 {
		req.Name = args[0]
	} else {
		fmt.Print("\tName: ")
		req.Name, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	ctx := withToken(context.Background())

	_, err = client.PasswordDelete(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Println("\tУдалён.")
	return nil
}

func updPasswordCmd(args []string) (err error) {
	var (
		readReq pb.PasswordReadRequest
		updReq  pb.PasswordUpdateRequest
	)

	if len(args) > 0 {
		readReq.Name = args[0]
	} else {
		fmt.Print("\t    Name: ")
		readReq.Name, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	ctx := withToken(context.Background())

	resp, err := client.PasswordRead(ctx, &readReq)
	if err != nil {
		return err
	}
	fmt.Println("\tUsername:", resp.Username)
	fmt.Println("\tPassword:", resp.Password)
	fmt.Println("\t   Notes:", resp.Notes)

	fmt.Println("\tNew")
	updReq.Id = resp.Id
	updReq.Write = &pb.PasswordWriteRequest{
		Name:     resp.Name,
		Username: resp.Username,
		Password: resp.Password,
		Notes:    resp.Notes,
	}

	fmt.Print("\t    Name: ")
	val, err := conReader.ReadLine()
	if err != nil {
		return err
	} else if val != "" {
		updReq.Write.Name = val
	}

	fmt.Print("\tUsername: ")
	if val, err = conReader.ReadLine(); err != nil {
		return err
	} else if val != "" {
		updReq.Write.Username = val
	}

	fmt.Print("\tPassword: ")
	if val, err = conReader.ReadLine(); err != nil {
		return err
	} else if val != "" {
		updReq.Write.Password = val
	}

	fmt.Print("\t   Notes: ")
	if val, err = conReader.ReadLine(); err != nil {
		return err
	} else if val != "" {
		updReq.Write.Notes = val
	}

	_, err = client.PasswordUpdate(ctx, &updReq)
	if err != nil {
		return err
	}
	fmt.Println("\tОбновлён.")
	return nil
}
