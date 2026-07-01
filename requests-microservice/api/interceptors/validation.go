package interceptors

import (
	"context"
	"requests-service/api/server"

	requestspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/requests/langs/go"
	"google.golang.org/grpc"
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
	case *requestspb.CreateRequestRequest:
		return server.ValidateCreateRequest(r)
	case *requestspb.GetRequestRequest:
		return server.ValidateGetRequest(r)
	case *requestspb.UpdateStatusRequest:
		return server.ValidateUpdateStatus(r)
	case *requestspb.AddCommentRequest:
		return server.ValidateAddComment(r)
	case *requestspb.GetCommentsRequest:
		return server.ValidateGetComments(r)
	}
	return nil
}
