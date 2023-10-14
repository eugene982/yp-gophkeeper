package grpc

import (
	"context"
	"time"

	"github.com/eugene982/yp-gophkeeper/internal/logger"
	"google.golang.org/grpc"
)

// loggerInterceptor прослойка логирования запросов
func loggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()

	logger.Info("request",
		"method", info.FullMethod)

	logger.Debug("incoming",
		"request", req)

	resp, err = handler(ctx, req)

	logger.Info("request",
		"duration", time.Since(start),
		"ok", err == nil)

	logger.Debug("outgoing",
		"response", resp,
		"error", err)

	return
}
