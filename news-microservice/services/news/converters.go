package news_service

import "news-service/infra/models/domain"

func createNewsRequestToModel(req *domain.CreateNewsRequest) *domain.News {
	return &domain.News{Title: req.Title, Content: req.Content, HouseID: req.HouseID, CreatedBy: req.CreatedBy}
}

func newsToCreatedEvent(news *domain.News) domain.NewsCreatedEvent {
	return domain.NewsCreatedEvent{Type: "news.created", NewsID: news.ID, Title: news.Title, HouseID: news.HouseID}
}
