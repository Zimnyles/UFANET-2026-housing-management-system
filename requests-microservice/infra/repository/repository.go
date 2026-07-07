package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	infra_errors "requests-service/infra/errors"
	"requests-service/infra/models/domain"
)

type Repository struct {
	db *gorm.DB
}

type dbMaintenanceRequest struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text;not null"`
	Type        string    `gorm:"type:varchar(32);not null;index"`
	Status      string    `gorm:"type:varchar(32);not null;index"`
	UserID      string    `gorm:"column:user_id;type:uuid;not null;index"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (dbMaintenanceRequest) TableName() string {
	return "maintenance_requests"
}

type dbComment struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RequestID string    `gorm:"column:request_id;type:uuid;not null;index"`
	UserID    string    `gorm:"column:user_id;type:uuid;not null;index"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (dbComment) TableName() string {
	return "request_comments"
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate(ctx context.Context) error {
	if err := r.db.WithContext(ctx).AutoMigrate(&dbMaintenanceRequest{}, &dbComment{}); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, req *domain.MaintenanceRequest) (*domain.MaintenanceRequest, error) {
	row := dbMaintenanceRequest{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Status:      domain.StatusOpen,
		UserID:      req.UserID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	return requestToDomain(&row), nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*domain.MaintenanceRequest, error) {
	var row dbMaintenanceRequest
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, infra_errors.ErrRequestNotFound
		}

		return nil, fmt.Errorf("get request: %w", err)
	}

	return requestToDomain(&row), nil
}

func (r *Repository) List(ctx context.Context, req *domain.GetRequestsRequest) ([]*domain.MaintenanceRequest, int64, error) {
	query := r.db.WithContext(ctx).Model(&dbMaintenanceRequest{})
	if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count requests: %w", err)
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	var rows []dbMaintenanceRequest
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list requests: %w", err)
	}

	result := make([]*domain.MaintenanceRequest, 0, len(rows))
	for i := range rows {
		result = append(result, requestToDomain(&rows[i]))
	}

	return result, total, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id string, status string) (*domain.MaintenanceRequest, error) {
	var row dbMaintenanceRequest
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, infra_errors.ErrRequestNotFound
		}

		return nil, fmt.Errorf("get request: %w", err)
	}

	row.Status = status

	row.UpdatedAt = time.Now().UTC()
	if err := r.db.WithContext(ctx).Save(&row).Error; err != nil {
		return nil, fmt.Errorf("update request status: %w", err)
	}

	return requestToDomain(&row), nil
}

func (r *Repository) AddComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	if _, err := r.GetByID(ctx, comment.RequestID); err != nil {
		return nil, err
	}

	row := dbComment{
		RequestID: comment.RequestID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: time.Now().UTC(),
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, fmt.Errorf("add comment: %w", err)
	}

	return commentToDomain(&row), nil
}

func (r *Repository) ListComments(ctx context.Context, requestID string) ([]*domain.Comment, error) {
	var rows []dbComment
	if err := r.db.WithContext(ctx).Where("request_id = ?", requestID).Order("created_at ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}

	result := make([]*domain.Comment, 0, len(rows))
	for i := range rows {
		result = append(result, commentToDomain(&rows[i]))
	}

	return result, nil
}

func requestToDomain(row *dbMaintenanceRequest) *domain.MaintenanceRequest {
	return &domain.MaintenanceRequest{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		Type:        row.Type,
		Status:      row.Status,
		UserID:      row.UserID,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func commentToDomain(row *dbComment) *domain.Comment {
	return &domain.Comment{
		ID:        row.ID,
		RequestID: row.RequestID,
		UserID:    row.UserID,
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
	}
}
