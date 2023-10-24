package main

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/eugene982/yp-gophkeeper/gen/go/proto/v1"
	"github.com/golang/protobuf/ptypes/empty"
)

func cardsCmd([]string) error {
	if userName == "" {
		return errUnauthenticated
	}
	currTable = "cards"
	return nil
}

func helpCardCmd(args []string) error {
	fmt.Println(`	
	"ls"         - список сохраненных карт
	"new [name]" - создание новой карты
	"get [name]" - получение из хранилища
	"upd [name]" - обновление данных в хранилище
	"del [name]" - получение пароля из хранилища`)
	return nil
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

func newCardCmd(args []string) (err error) {
	var req pb.CardWriteRequest

	if len(args) > 0 {
		req.Name = args[0]
	} else {
		fmt.Print("\t  Name: ")
		req.Name, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	fmt.Print("\tNumber: ")
	req.Number, err = conReader.ReadLine()
	if err != nil {
		return err
	}

	fmt.Print("\t   Pin: ")
	req.Pin, err = conReader.ReadLine()
	if err != nil {
		return err
	}

	fmt.Print("\t Notes: ")
	req.Notes, err = conReader.ReadLine()
	if err != nil {
		return err
	}

	ctx := withToken(context.Background())
	_, err = client.CardWrite(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Println("\tДобавлено.")
	return nil
}

func getCardCmd(args []string) (err error) {
	var req pb.CardReadRequest

	if len(args) > 0 {
		req.Name = args[0]
	} else {
		fmt.Print("\t  Name: ")
		req.Name, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	ctx := withToken(context.Background())

	resp, err := client.CardRead(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Println("\tNumber:", resp.Number)
	fmt.Println("\t   Pin:", resp.Pin)
	fmt.Println("\t Notes:", resp.Notes)
	return nil
}

func delCardCmd(args []string) (err error) {
	var req pb.CardDelRequest

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

	_, err = client.CardDelete(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Println("\tУдалено.")
	return nil
}

func updCardCmd(args []string) (err error) {
	var (
		readReq pb.CardReadRequest
		updReq  pb.CardUpdateRequest
	)

	if len(args) > 0 {
		readReq.Name = args[0]
	} else {
		fmt.Print("\t  Name: ")
		readReq.Name, err = conReader.ReadLine()
		if err != nil {
			return err
		}
	}

	ctx := withToken(context.Background())

	resp, err := client.CardRead(ctx, &readReq)
	if err != nil {
		return err
	}
	fmt.Println("\tNumber:", resp.Number)
	fmt.Println("\t   Pin:", resp.Pin)
	fmt.Println("\t Notes:", resp.Notes)

	fmt.Println("\tNew")
	updReq.Id = resp.Id
	updReq.Write = &pb.CardWriteRequest{
		Name:   resp.Name,
		Number: resp.Number,
		Pin:    resp.Pin,
		Notes:  resp.Notes,
	}

	fmt.Print("\t  Name: ")
	val, err := conReader.ReadLine()
	if err != nil {
		return err
	} else if val != "" {
		updReq.Write.Name = val
	}

	fmt.Print("\t Number: ")
	if val, err = conReader.ReadLine(); err != nil {
		return err
	} else if val != "" {
		updReq.Write.Number = val
	}

	fmt.Print("\t   Pin: ")
	if val, err = conReader.ReadLine(); err != nil {
		return err
	} else if val != "" {
		updReq.Write.Pin = val
	}

	fmt.Print("\t  Notes: ")
	if val, err = conReader.ReadLine(); err != nil {
		return err
	} else if val != "" {
		updReq.Write.Notes = val
	}

	_, err = client.CardUpdate(ctx, &updReq)
	if err != nil {
		return err
	}
	fmt.Println("\tОбновлено.")
	return nil
}
