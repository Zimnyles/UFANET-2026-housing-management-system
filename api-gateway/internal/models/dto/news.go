package dto

type CreateNewsRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	HouseID string `json:"house_id"`
}

type NewsResponse struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	HouseID    string `json:"house_id"`
	CreatedAt  string `json:"created_at"`
	AuthorName string `json:"author_name"`
}

type ListNewsResponse struct {
	News  []NewsResponse `json:"news"`
	Total int64          `json:"total"`
}
