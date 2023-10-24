package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/golang/protobuf/ptypes/empty"
)

func filesCmd([]string) error {
	if userName == "" {
		return errUnauthenticated
	}
	currTable = "files"
	return nil
}

func helpFileCmd(args []string) error {
	fmt.Println(`	
	"ls"         - список сохраненных файлов
	"new [name]" - создание нового файла
	"get [name]" - получение из хранилища
	"upd [name]" - обновление данных в хранилище
	"del [name]" - получение пароля из хранилища`)
	return nil
}

func lsFileCmd([]string) error {

	ctx := withToken(context.Background())
	resp, err := client.BinaryList(ctx, &empty.Empty{})
	if err != nil {
		return nil
	}

	fmt.Println("\tFiles:", strings.Join(resp.Names, "; "))
	fmt.Println("\tCount:", len(resp.Names))

	return nil
}

func newFileCmd(args []string) (err error) {
	var req pb.BinaryWriteRequest

	var filename string
	if len(args) > 0 {
		filename = args[0]
	} else {
		fmt.Print("\tFilename: ")
		filename, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	req.Name = path.Base(filename)
	req.Bin, err = readFile(filename)
	if err != nil {
		return err
	}

	fmt.Print("\t   Notes: ")
	req.Notes, err = conReader.ReadLine()
	if err != nil {
		return err
	}

	ctx := withToken(context.Background())
	_, err = client.BinaryWrite(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Println("\tДобавлено.")
	return nil
}

func getFileCmd(args []string) (err error) {
	var req pb.BinaryReadRequest

	if len(args) > 0 {
		req.Name = args[0]
	} else {
		fmt.Print("\t Name: ")
		req.Name, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	ctx := withToken(context.Background())

	resp, err := client.BinaryRead(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Println("\tNotes:", resp.Notes)

	fmt.Print("\tFilename: ")
	filename, err := conReader.ReadLine()
	if err != nil {
		return err
	} else if filename == "" {
		return nil
	}
	return writeFile(filename, resp.Bin)
}

func delFileCmd(args []string) (err error) {
	var req pb.BinaryDelRequest

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

	_, err = client.BinaryDelete(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Println("\tУдалено.")
	return nil
}

func updFileCmd(args []string) (err error) {
	var (
		readReq pb.BinaryReadRequest
		updReq  pb.BinaryUpdateRequest
	)

	if len(args) > 0 {
		readReq.Name = args[0]
	} else {
		fmt.Print("\t   Name: ")
		readReq.Name, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	ctx := withToken(context.Background())

	resp, err := client.BinaryRead(ctx, &readReq)
	if err != nil {
		return err
	}
	fmt.Println("\t  Name:", resp.Name)
	fmt.Println("\t Notes:", resp.Notes)

	fmt.Println("\tNew")
	updReq.Id = resp.Id
	updReq.Write = &pb.BinaryWriteRequest{
		Name:  resp.Name,
		Notes: resp.Notes,
	}

	fmt.Print("\t  Name: ")
	val, err := conReader.ReadLine()
	if err != nil {
		return err
	} else if val != "" {
		updReq.Write.Name = val
	}

	fmt.Print("\tNotes: ")
	if val, err = conReader.ReadLine(); err != nil {
		return err
	} else if val != "" {
		updReq.Write.Notes = val
	}

	fmt.Print("\tFilename: ")
	if val, err = conReader.ReadLine(); err != nil {
		return err
	} else if val != "" {
		updReq.Write.Bin, err = readFile(val)
		if err != nil {
			return err
		}
	}

	_, err = client.BinaryUpdate(ctx, &updReq)
	if err != nil {
		return err
	}
	fmt.Println("\tОбновлено.")
	return nil
}

func readFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

func writeFile(fileName string, b []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(b)
	return err
}
