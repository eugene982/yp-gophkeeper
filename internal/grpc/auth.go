package grpc

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/eugene982/yp-gophkeeper/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	TOKEN_EXP  = time.Hour * 3
	SECRET_KEY = "supersecretkey"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

// SetUserID устанавливает идентификатор пользователя в контекст запроса
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextKeyUserID, userID)
}

// GetUserID Возвращает идентификатор пользователя из контекста
func GetUserID(ctx context.Context) (string, error) {
	val := ctx.Value(contextKeyUserID)
	if val == nil {
		return "", fmt.Errorf("user id not found")
	}
	userID, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("user id is not uint type")
	}
	return userID, nil
}

func newAuthInterceptor(handlers ...string) grpc.UnaryServerInterceptor {

	handlersSlice := sort.StringSlice(handlers)
	handlersSlice.Sort()

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		method, ok := strings.CutSuffix(info.FullMethod, "/")
		if !ok || handlersSlice.Search(method) == -1 {
			return handler(ctx, req)
		}

		var token string
		if md, ok := metadata.FromIncomingContext(ctx); !ok {
			return nil, errMissingMetadata
		} else if vals := md.Get("token"); len(vals) > 1 {
			token = vals[0]
		} else {
			return nil, errInvalidToken
		}

		userID, err := utils.GetJWTUserID(token, SECRET_KEY)
		if err != nil {
			return nil, errInvalidToken
		}

		// помещаем идентификатор пользователя в контекст и выполняем с ним метод
		uctx := SetUserID(ctx, userID)
		return handler(uctx, req)
	}
}
