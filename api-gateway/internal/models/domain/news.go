package domain

type News struct {
	ID        string
	Title     string
	Content   string
	HouseID   string
	CreatedAt string
	CreatedBy string
}

type CreateNews struct {
	Title     string
	Content   string
	HouseID   string
	CreatedBy string
}

type ListNews struct {
	HouseID  string
	DateFrom string
	DateTo   string
	Limit    int
	Offset   int
}
