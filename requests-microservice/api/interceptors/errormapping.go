package interceptors

import (
	"context"

	"google.golang.org/grpc"

	"requests-service/infra/grpcerrors"
)

func ErrorMapping() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, grpcerrors.ToGrpcError(err)
		}

		return resp, nil
	}
}
