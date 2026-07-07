package news_handler

import (
	"api-gateway/internal/models/domain"
	"api-gateway/internal/models/dto"
)

func toDomainCreate(req dto.CreateNewsRequest, houseID, userID string) *domain.CreateNews {
	return &domain.CreateNews{Title: req.Title, Content: req.Content, HouseID: houseID, CreatedBy: userID}
}

func toDomainList(houseID, from, to string, limit, offset int) *domain.ListNews {
	return &domain.ListNews{HouseID: houseID, DateFrom: from, DateTo: to, Limit: limit, Offset: offset}
}

func toDTO(news *domain.News, authorName string) dto.NewsResponse {
	return dto.NewsResponse{ID: news.ID, Title: news.Title, Content: news.Content, HouseID: news.HouseID, CreatedAt: news.CreatedAt, AuthorName: authorName}
}
