package domain

import "time"

type Device struct{ UserID, Token, Platform string }

type Notification struct {
	ID, UserID, HouseID, Type, Title, Body string
	CreatedAt                              time.Time
	Read                                   bool
}

type ListRequest struct {
	UserID, HouseID string
	Limit, Offset   int
}
