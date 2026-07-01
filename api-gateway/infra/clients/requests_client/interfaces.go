package requests_client

import (
	"context"

	requestspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/requests/langs/go"
	"google.golang.org/grpc"
)

type ClientConnection interface {
	Close() error
}

type RequestsServiceClient interface {
	CreateRequest(ctx context.Context, in *requestspb.CreateRequestRequest, opts ...grpc.CallOption) (*requestspb.RequestResponse, error)
	GetRequests(ctx context.Context, in *requestspb.GetRequestsRequest, opts ...grpc.CallOption) (*requestspb.GetRequestsResponse, error)
	GetRequest(ctx context.Context, in *requestspb.GetRequestRequest, opts ...grpc.CallOption) (*requestspb.RequestResponse, error)
	UpdateRequestStatus(ctx context.Context, in *requestspb.UpdateStatusRequest, opts ...grpc.CallOption) (*requestspb.RequestResponse, error)
	AddComment(ctx context.Context, in *requestspb.AddCommentRequest, opts ...grpc.CallOption) (*requestspb.CommentResponse, error)
	GetComments(ctx context.Context, in *requestspb.GetCommentsRequest, opts ...grpc.CallOption) (*requestspb.GetCommentsResponse, error)
}
