package news_client

import (
	"math"

	newspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/news/langs/go"

	"api-gateway/internal/models/domain"
)

func toProtoCreate(req *domain.CreateNews) *newspb.CreateNewsRequest {
	return &newspb.CreateNewsRequest{Title: req.Title, Content: req.Content, HouseId: req.HouseID, CreatedBy: req.CreatedBy}
}

func toProtoList(req *domain.ListNews) *newspb.GetNewsRequest {
	return &newspb.GetNewsRequest{
		HouseId: req.HouseID, DateFrom: req.DateFrom, DateTo: req.DateTo,
		Limit: paginationInt32(req.Limit), Offset: paginationInt32(req.Offset),
	}
}

func paginationInt32(value int) int32 {
	if value <= 0 {
		return 0
	}

	if value > math.MaxInt32 {
		return math.MaxInt32
	}

	return int32(value)
}

func toDomain(pb *newspb.News) *domain.News {
	if pb == nil {
		return nil
	}

	return &domain.News{ID: pb.GetId(), Title: pb.GetTitle(), Content: pb.GetContent(), HouseID: pb.GetHouseId(), CreatedAt: pb.GetCreatedAt(), CreatedBy: pb.GetCreatedBy()}
}

func toDomainList(items []*newspb.News) []*domain.News {
	result := make([]*domain.News, 0, len(items))
	for _, item := range items {
		result = append(result, toDomain(item))
	}

	return result
}
