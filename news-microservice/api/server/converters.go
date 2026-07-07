package server

import (
	"time"

	newspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/news/langs/go"

	"news-service/infra/models/domain"
)

func protoToCreateNews(req *newspb.CreateNewsRequest) *domain.CreateNewsRequest {
	return &domain.CreateNewsRequest{
		Title:     req.GetTitle(),
		Content:   req.GetContent(),
		HouseID:   req.GetHouseId(),
		CreatedBy: req.GetCreatedBy(),
	}
}

func protoToGetNews(req *newspb.GetNewsRequest) *domain.GetNewsRequest {
	out := &domain.GetNewsRequest{
		HouseID: req.GetHouseId(),
		Limit:   int(req.GetLimit()),
		Offset:  int(req.GetOffset()),
	}
	if from := req.GetDateFrom(); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			out.DateFrom = t
		}
	}

	if to := req.GetDateTo(); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			out.DateTo = t
		}
	}

	return out
}

func domainToProtoNews(n *domain.News) *newspb.News {
	return &newspb.News{
		Id:        n.ID,
		Title:     n.Title,
		Content:   n.Content,
		HouseId:   n.HouseID,
		CreatedAt: n.CreatedAt.UTC().Format(time.RFC3339),
		CreatedBy: n.CreatedBy,
	}
}
