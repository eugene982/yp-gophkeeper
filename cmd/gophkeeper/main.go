package main

import (
	"log"

	"github.com/eugene982/yp-gophkeeper/internal/config"
	"github.com/eugene982/yp-gophkeeper/internal/logger"
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

	return nil
}
