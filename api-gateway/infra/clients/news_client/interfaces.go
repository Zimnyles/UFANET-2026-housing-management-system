package news_client

import (
	"context"

	newspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/news/langs/go"
	"google.golang.org/grpc"
)

type ClientConnection interface{ Close() error }

type NewsServiceClient interface {
	CreateNews(context.Context, *newspb.CreateNewsRequest, ...grpc.CallOption) (*newspb.NewsResponse, error)
	GetNews(context.Context, *newspb.GetNewsRequest, ...grpc.CallOption) (*newspb.GetNewsResponse, error)
}
