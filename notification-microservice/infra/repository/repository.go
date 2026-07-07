package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"notification-service/infra/models/domain"
)

type dbDevice struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      string `gorm:"type:uuid;not null;uniqueIndex:idx_user_device"`
	DeviceToken string `gorm:"not null;uniqueIndex:idx_user_device"`
	Platform    string `gorm:"not null"`
	CreatedAt   time.Time
}

func (dbDevice) TableName() string { return "notification_devices" }

type dbNotification struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    *string   `gorm:"type:uuid;index"`
	HouseID   *string   `gorm:"type:uuid;index"`
	Type      string    `gorm:"not null;index"`
	Title     string    `gorm:"not null"`
	Body      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;index"`
	Read      bool      `gorm:"not null;default:false"`
}

func (dbNotification) TableName() string { return "notifications" }

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }
func (r *Repository) Migrate(ctx context.Context) error {
	return r.db.WithContext(ctx).AutoMigrate(&dbDevice{}, &dbNotification{})
}

func (r *Repository) Register(ctx context.Context, device *domain.Device) error {
	row := deviceToRow(device)

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "user_id"}, {Name: "device_token"}}, DoUpdates: clause.AssignmentColumns([]string{"platform"})}).Create(row).Error
}

func (r *Repository) Unregister(ctx context.Context, userID, token string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND device_token = ?", userID, token).Delete(&dbDevice{}).Error
}

func (r *Repository) Create(ctx context.Context, n *domain.Notification) (*domain.Notification, error) {
	row := notificationToRow(n)
	if row.CreatedAt.IsZero() {
		row.CreatedAt = time.Now().UTC()
	}

	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}

	return notificationToDomain(row), nil
}

func (r *Repository) List(ctx context.Context, req *domain.ListRequest) ([]*domain.Notification, int64, error) {
	query := r.db.WithContext(ctx).Model(&dbNotification{}).Where("user_id = ?", req.UserID)
	if req.HouseID != "" {
		query = r.db.WithContext(ctx).Model(&dbNotification{}).Where("user_id = ? OR house_id = ?", req.UserID, req.HouseID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	if req.Offset < 0 {
		req.Offset = 0
	}

	var rows []dbNotification
	if err := query.Order("created_at DESC").Limit(limit).Offset(req.Offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*domain.Notification, 0, len(rows))
	for i := range rows {
		result = append(result, notificationToDomain(&rows[i]))
	}

	return result, total, nil
}
