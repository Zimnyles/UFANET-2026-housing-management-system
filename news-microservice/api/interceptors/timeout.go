package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func Timeout(d time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, d)
		defer cancel()

		return handler(ctx, req)
	}
}
