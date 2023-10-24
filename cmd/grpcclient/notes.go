package main

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/golang/protobuf/ptypes/empty"
)

func notesCmd([]string) error {
	if userName == "" {
		return errUnauthenticated
	}
	currTable = "notes"
	return nil
}

func helpNoteCmd(args []string) error {
	fmt.Println(`	
	"ls"         - список сохраненных заметок
	"new [name]" - создание новой заметки
	"get [name]" - получение из хранилища
	"upd [name]" - обновление данных в хранилище
	"del [name]" - получение пароля из хранилища`)
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

func newNoteCmd(args []string) (err error) {
	var req pb.NoteWriteRequest

	if len(args) > 0 {
		req.Name = args[0]
	} else {
		fmt.Print("\t Name: ")
		req.Name, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	fmt.Print("\tNotes: ")
	req.Notes, err = conReader.ReadLine()
	if err != nil {
		return err
	}

	ctx := withToken(context.Background())
	_, err = client.NoteWrite(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Println("\tДобавлено.")
	return nil
}

func getNoteCmd(args []string) (err error) {
	var req pb.NoteReadRequest

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

	resp, err := client.NoteRead(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Println("\tNotes:", resp.Notes)
	return nil
}

func delNoteCmd(args []string) (err error) {
	var req pb.NoteDelRequest

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

	_, err = client.NoteDelete(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Println("\tУдалено.")
	return nil
}

func updNoteCmd(args []string) (err error) {
	var (
		readReq pb.NoteReadRequest
		updReq  pb.NoteUpdateRequest
	)

	if len(args) > 0 {
		readReq.Name = args[0]
	} else {
		fmt.Print("\t Name: ")
		readReq.Name, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	ctx := withToken(context.Background())

	resp, err := client.NoteRead(ctx, &readReq)
	if err != nil {
		return err
	}
	fmt.Println("\tNotes:", resp.Notes)

	fmt.Println("\tNew")
	updReq.Id = resp.Id
	updReq.Write = &pb.NoteWriteRequest{
		Name:  resp.Name,
		Notes: resp.Notes,
	}

	fmt.Print("\t Name: ")
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

	_, err = client.NoteUpdate(ctx, &updReq)
	if err != nil {
		return err
	}
	fmt.Println("\tОбновлено.")
	return nil
}
