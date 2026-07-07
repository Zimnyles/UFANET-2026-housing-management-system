package interceptors

import (
	"context"

	newspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/news/langs/go"
	"google.golang.org/grpc"

	"news-service/api/server"
)

func Validation() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := validate(req); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func validate(req interface{}) error {
	switch r := req.(type) {
	case *newspb.CreateNewsRequest:
		return server.ValidateCreateNewsRequest(r)
	case *newspb.GetNewsItemRequest:
		return server.ValidateGetNewsItemRequest(r)
	}

	return nil
}
