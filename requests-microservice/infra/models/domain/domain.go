package domain

import "time"

const (
	RequestTypePlumber     = "plumber"
	RequestTypeElectrician = "electrician"

	StatusOpen       = "open"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
	StatusCancelled  = "cancelled"
)

type MaintenanceRequest struct {
	ID          string
	Title       string
	Description string
	Type        string
	Status      string
	UserID      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Comment struct {
	ID        string
	RequestID string
	UserID    string
	Content   string
	CreatedAt time.Time
}

type CreateRequestRequest struct {
	Title       string
	Description string
	Type        string
	UserID      string
}

type GetRequestsRequest struct {
	UserID string
	Status string
	Type   string
	Limit  int
	Offset int
}

type UpdateStatusRequest struct {
	ID     string
	Status string
	UserID string
}

type AddCommentRequest struct {
	RequestID string
	UserID    string
	Content   string
}

type RequestEvent struct {
	Type      string `json:"type"`
	RequestID string `json:"request_id"`
	UserID    string `json:"user_id,omitempty"`
	Status    string `json:"status,omitempty"`
	Title     string `json:"title,omitempty"`
}
