package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	infra_errors "auth-service/infra/errors"
	"auth-service/infra/models/domain"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type AdminCode struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code      string    `gorm:"uniqueIndex;not null"`
	Status    string    `gorm:"not null;default:active"`
	CreatedAt time.Time `gorm:"not null"`
}

func (AdminCode) TableName() string {
	return "admin_codes"
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate(ctx context.Context) error {
	db := r.db.WithContext(ctx)
	if err := db.AutoMigrate(&domain.User{}, &AdminCode{}); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	if err := db.Where(AdminCode{Code: "123123123"}).FirstOrCreate(&AdminCode{
		Code:   "123123123",
		Status: "active",
	}).Error; err != nil {
		return fmt.Errorf("seed admin code: %w", err)
	}
	return nil
}

func (r *Repository) IsActiveAdminCode(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&AdminCode{}).
		Where("code = ? AND status = ?", code, "active").
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("check admin code: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || containsCode(err, "23505") {
			return nil, infra_errors.ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, infra_errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil
}

func containsCode(err error, code string) bool {
	type pgErr interface {
		SQLState() string
	}
	if e, ok := err.(pgErr); ok {
		return e.SQLState() == code
	}
	return false
}
