package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v8"
)

// Config конфигурация получаемая из флагов и/или переменных окружения
type Config struct {
	ServerAddres string        `env:"RUN_ADDRESS"`
	LogLevel     string        `env:"LOG_LEVEL"` // уровень логирования
	Timeout      time.Duration // таймаут соединения
}

// Parse заполнение структуры конфигурации
func Parse() (Config, error) {
	var config Config

	// читаем флаги
	flag.StringVar(&config.ServerAddres, "a", ":28000", "server address")
	flag.DurationVar(&config.Timeout, "t", 30, "timeout in seconds")
	flag.StringVar(&config.LogLevel, "l", "info", "log level")
	flag.Parse()

	err := env.Parse(&config)
	return config, err
}
