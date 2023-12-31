// Package config конфигурация читаемая из флагов и переменных окружения
package config

import (
	"flag"

	"github.com/caarlos0/env/v8"
)

// Config конфигурация получаемая из флагов и/или переменных окружения
type Config struct {
	ServerAddres string `env:"RUN_ADDRESS"`
	LogLevel     string `env:"LOG_LEVEL"`    // уровень логирования
	DSN          string `env:"DATABASE_DSN"` // адрес подключения к базе данных
	MigratePath  string `env:"MIGRATE_PATH"` // адрес подключения к базе данных
}

// Parse заполнение структуры конфигурации
func Parse() (Config, error) {
	var config Config

	// читаем флаги
	flag.StringVar(&config.ServerAddres, "a", ":28000", "server address")
	flag.StringVar(&config.LogLevel, "l", "info", "log level")
	flag.StringVar(&config.DSN, "d", "postgres://postgres:postgres@localhost/gophkeeper", "postgres connection string")
	flag.StringVar(&config.MigratePath, "m", "", "path to migrations files")
	flag.Parse()

	err := env.Parse(&config)
	return config, err
}
