package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	infra_errors "news-service/infra/errors"
	"news-service/infra/models/domain"
)

type Repository struct {
	db *gorm.DB
}

type dbNews struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title     string    `gorm:"type:varchar(255);not null"`
	Content   string    `gorm:"type:text;not null"`
	HouseID   string    `gorm:"column:house_id;type:uuid;not null;index"`
	CreatedAt time.Time `gorm:"not null"`
	CreatedBy string    `gorm:"column:created_by;type:uuid;index"`
}

func (dbNews) TableName() string { return "news" }

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate(ctx context.Context) error {
	if err := r.db.WithContext(ctx).AutoMigrate(&dbNews{}); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, n *domain.News) (*domain.News, error) {
	row := dbNews{
		Title:     n.Title,
		Content:   n.Content,
		HouseID:   n.HouseID,
		CreatedAt: time.Now().UTC(),
		CreatedBy: n.CreatedBy,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, fmt.Errorf("create news: %w", err)
	}

	return newsToDomain(&row), nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*domain.News, error) {
	var row dbNews
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, infra_errors.ErrNewsNotFound
		}

		return nil, fmt.Errorf("get news: %w", err)
	}

	return newsToDomain(&row), nil
}

func (r *Repository) List(ctx context.Context, req *domain.GetNewsRequest) ([]*domain.News, int64, error) {
	query := r.db.WithContext(ctx).Model(&dbNews{})

	if req.HouseID != "" {
		query = query.Where("house_id = ?", req.HouseID)
	}

	if !req.DateFrom.IsZero() {
		query = query.Where("created_at >= ?", req.DateFrom)
	}

	if !req.DateTo.IsZero() {
		query = query.Where("created_at <= ?", req.DateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count news: %w", err)
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	var rows []dbNews
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list news: %w", err)
	}

	result := make([]*domain.News, 0, len(rows))
	for i := range rows {
		result = append(result, newsToDomain(&rows[i]))
	}

	return result, total, nil
}

func newsToDomain(row *dbNews) *domain.News {
	return &domain.News{
		ID:        row.ID,
		Title:     row.Title,
		Content:   row.Content,
		HouseID:   row.HouseID,
		CreatedAt: row.CreatedAt,
		CreatedBy: row.CreatedBy,
	}
}
