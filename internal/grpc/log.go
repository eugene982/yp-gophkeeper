package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/eugene982/yp-gophkeeper/internal/logger"
)

// loggerInterceptor прослойка логирования запросов
func loggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()

	logger.Info("request",
		"method", info.FullMethod)

	md, _ := metadata.FromIncomingContext(ctx)

	logger.Debug("incoming",
		"request", req,
		"metadata", md)

	resp, err = handler(ctx, req)

	logger.Info("response",
		"duration", time.Since(start),
		"ok", err == nil)

	logger.Debug("outgoing",
		"response", resp,
		"error", err)

	return
}
