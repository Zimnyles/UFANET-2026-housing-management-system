package dto

type MaintenanceRequestResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	UserID      string `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	AuthorName  string `json:"author_name"`
}

type RequestCommentResponse struct {
	ID         string `json:"id"`
	RequestID  string `json:"request_id"`
	UserID     string `json:"user_id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	AuthorName string `json:"author_name"`
}

type CreateMaintenanceRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type ListMaintenanceRequestsResponse struct {
	Requests []MaintenanceRequestResponse `json:"requests"`
	Total    int64                        `json:"total"`
}

type UpdateMaintenanceRequestStatus struct {
	Status string `json:"status"`
}

type AddMaintenanceRequestComment struct {
	Content string `json:"content"`
}

type ListMaintenanceRequestCommentsResponse struct {
	Comments []RequestCommentResponse `json:"comments"`
}
