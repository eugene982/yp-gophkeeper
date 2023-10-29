package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/eugene982/yp-gophkeeper/internal/application"
	"github.com/eugene982/yp-gophkeeper/internal/config"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
)

const (
	// сколько ждём времени на корректное завершение работы сервера
	closeServerTimeout = time.Second * 3
)

var (
	buildVersion, buildDate, buildCommit string
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	// логируем любую возврвщаемую ошибку
	defer func() {
		if err != nil {
			logger.Errorf("error stating server: %w", err)
		}
	}()

	config, err := config.Parse()
	if err != nil {
		return
	}

	err = logger.Initialize(config.LogLevel)
	if err != nil {
		return
	}
	logger.Debug("debug enable")
	logger.Info("build",
		"version", buildVersion,
		"date", buildDate,
		"commit", buildCommit)

	app, err := application.New(config)
	if err != nil {
		return
	}

	// захват прерывания процесса
	ctxInterrupt, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// запуск сервера в горутине
	srvErr := make(chan error)
	go func() {
		srvErr <- app.Start()
	}()
	logger.Info("application start", "config", config)

	// ждём что раньше случится, ошибка старта сервера
	// или пользователь прервёт программу
	select {
	case <-ctxInterrupt.Done():
		// прервано пользователем
	case e := <-srvErr:
		// сервер не смог стартануть, некорректый адрес, занят порт...
		// эту ошибку логируем отдельно. В любом случае, нужно освободить ресурсы
		logger.Errorf("error start server: %w", e)
	}

	// стартуем завершение сервера
	stopErr := make(chan error)
	go func() {
		stopErr <- app.Stop()
	}()

	// Ждём пока сервер сам завершится
	// или за отведённое время
	ctxTimeout, stop := context.WithTimeout(context.Background(), closeServerTimeout)
	defer stop()

	select {
	case <-ctxTimeout.Done():
		logger.Warn("stop server on timeout")
		return nil
	case err := <-stopErr:
		logger.Info("stop server gracefull")
		return err
	}
}
