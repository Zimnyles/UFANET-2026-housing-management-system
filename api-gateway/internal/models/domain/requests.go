package domain

type MaintenanceRequest struct {
	ID          string
	Title       string
	Description string
	Type        string
	Status      string
	UserID      string
	CreatedAt   string
	UpdatedAt   string
}

type RequestComment struct {
	ID        string
	RequestID string
	UserID    string
	Content   string
	CreatedAt string
}

type CreateMaintenanceRequest struct {
	Title       string
	Description string
	Type        string
	UserID      string
}

type ListMaintenanceRequests struct {
	UserID string
	Status string
	Type   string
	Limit  int
	Offset int
}

type UpdateMaintenanceRequestStatus struct {
	ID     string
	Status string
	UserID string
}

type AddMaintenanceRequestComment struct {
	RequestID string
	UserID    string
	Content   string
}
