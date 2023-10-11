package handler

import "context"

// Pinger интерфейс проверки соединения
type Pinger interface {
	Ping(context.Context) error
}

// Тип
type PingerFunc func(context.Context) error

func (f PingerFunc) Ping(ctx context.Context) error {
	return f(ctx)
}

var _ Pinger = PingerFunc(nil)
