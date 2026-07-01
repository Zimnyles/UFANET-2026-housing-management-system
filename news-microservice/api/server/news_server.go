package server

import (
	"context"

	newspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/news/langs/go"
	"github.com/rs/zerolog"
)

type NewsServer struct {
	newspb.UnimplementedNewsServiceServer
	service NewsService
	logger  *zerolog.Logger
}

func NewNewsServer(service NewsService, logger *zerolog.Logger) *NewsServer {
	return &NewsServer{service: service, logger: logger}
}

func (s *NewsServer) CreateNews(ctx context.Context, req *newspb.CreateNewsRequest) (*newspb.NewsResponse, error) {
	n, err := s.service.CreateNews(ctx, protoToCreateNews(req))
	if err != nil {
		return nil, err
	}
	return &newspb.NewsResponse{News: domainToProtoNews(n)}, nil
}

func (s *NewsServer) GetNewsItem(ctx context.Context, req *newspb.GetNewsItemRequest) (*newspb.NewsResponse, error) {
	n, err := s.service.GetNewsItem(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &newspb.NewsResponse{News: domainToProtoNews(n)}, nil
}

func (s *NewsServer) GetNews(ctx context.Context, req *newspb.GetNewsRequest) (*newspb.GetNewsResponse, error) {
	items, total, err := s.service.GetNews(ctx, protoToGetNews(req))
	if err != nil {
		return nil, err
	}
	pb := make([]*newspb.News, 0, len(items))
	for _, n := range items {
		pb = append(pb, domainToProtoNews(n))
	}
	return &newspb.GetNewsResponse{News: pb, Total: int32(total)}, nil
}
