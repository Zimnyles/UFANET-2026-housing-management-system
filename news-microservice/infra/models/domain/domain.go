package domain

import "time"

type News struct {
	ID        string
	Title     string
	Content   string
	HouseID   string
	CreatedAt time.Time
	CreatedBy string
}

type CreateNewsRequest struct {
	Title     string
	Content   string
	HouseID   string
	CreatedBy string
}

type GetNewsRequest struct {
	HouseID  string
	DateFrom time.Time
	DateTo   time.Time
	Limit    int
	Offset   int
}

type NewsCreatedEvent struct {
	Type    string `json:"type"`
	NewsID  string `json:"news_id"`
	Title   string `json:"title"`
	HouseID string `json:"house_id"`
}
