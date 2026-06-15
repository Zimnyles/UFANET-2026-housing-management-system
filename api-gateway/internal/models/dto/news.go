package dto

type CreateNewsRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
