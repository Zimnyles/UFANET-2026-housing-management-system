package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	infra_errors "profile-service/infra/errors"
	"profile-service/infra/models/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

type dbManagementCompany struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (dbManagementCompany) TableName() string {
	return "management_companies"
}

type dbHouse struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"not null"`
	Address   string    `gorm:"not null"`
	UKID      string    `gorm:"column:uk_id;type:uuid;not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (dbHouse) TableName() string {
	return "houses"
}

type dbProfile struct {
	UserID    string    `gorm:"column:user_id;type:uuid;primaryKey"`
	FullName  string    `gorm:"not null;default:''"`
	Phone     string    `gorm:"not null;default:''"`
	Apartment string    `gorm:"not null;default:''"`
	HouseID   string    `gorm:"column:house_id;type:uuid"`
	UKID      string    `gorm:"column:uk_id;type:uuid"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (dbProfile) TableName() string {
	return "profiles"
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate(ctx context.Context) error {
	if err := r.db.WithContext(ctx).AutoMigrate(&dbManagementCompany{}, &dbHouse{}, &dbProfile{}); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}

func (r *Repository) GetProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	var p dbProfile
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, infra_errors.ErrProfileNotFound
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return profileToDomain(&p), nil
}

func (r *Repository) UpsertProfile(ctx context.Context, profile *domain.Profile) (*domain.Profile, error) {
	db := r.db.WithContext(ctx)
	var house dbHouse
	if err := db.Where("id = ?", profile.HouseID).First(&house).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, infra_errors.ErrHouseIDInvalid
		}
		return nil, fmt.Errorf("get house: %w", err)
	}

	p := dbProfile{
		UserID:    profile.UserID,
		FullName:  profile.FullName,
		Phone:     profile.Phone,
		Apartment: profile.Apartment,
		HouseID:   house.ID,
		UKID:      house.UKID,
	}

	if err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"full_name",
			"phone",
			"apartment",
			"house_id",
			"uk_id",
			"updated_at",
		}),
	}).Create(&p).Error; err != nil {
		return nil, fmt.Errorf("upsert profile: %w", err)
	}

	return r.GetProfile(ctx, profile.UserID)
}

func (r *Repository) IsProfileComplete(ctx context.Context, userID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&dbProfile{}).
		Where("user_id = ? AND full_name <> '' AND phone <> '' AND apartment <> '' AND house_id <> ''", userID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("is profile complete: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) CreateManagementCompany(ctx context.Context, company *domain.ManagementCompany) (*domain.ManagementCompany, error) {
	c := dbManagementCompany{Name: company.Name}
	if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
		return nil, fmt.Errorf("create management company: %w", err)
	}
	return companyToDomain(&c), nil
}

func (r *Repository) ListManagementCompanies(ctx context.Context) ([]*domain.ManagementCompany, error) {
	var companies []dbManagementCompany
	if err := r.db.WithContext(ctx).Order("name").Find(&companies).Error; err != nil {
		return nil, fmt.Errorf("list management companies: %w", err)
	}
	result := make([]*domain.ManagementCompany, 0, len(companies))
	for i := range companies {
		result = append(result, companyToDomain(&companies[i]))
	}
	return result, nil
}

func (r *Repository) CreateHouse(ctx context.Context, house *domain.House) (*domain.House, error) {
	h := dbHouse{Name: house.Name, Address: house.Address, UKID: house.UKID}
	if err := r.db.WithContext(ctx).Create(&h).Error; err != nil {
		return nil, fmt.Errorf("create house: %w", err)
	}
	return houseToDomain(&h), nil
}

func (r *Repository) ListHouses(ctx context.Context, ukID string) ([]*domain.House, error) {
	var houses []dbHouse
	query := r.db.WithContext(ctx)
	if ukID != "" {
		query = query.Where("uk_id = ?", ukID)
	}
	if err := query.Order("name").Find(&houses).Error; err != nil {
		return nil, fmt.Errorf("list houses: %w", err)
	}
	result := make([]*domain.House, 0, len(houses))
	for i := range houses {
		result = append(result, houseToDomain(&houses[i]))
	}
	return result, nil
}

func profileToDomain(p *dbProfile) *domain.Profile {
	return &domain.Profile{
		UserID:    p.UserID,
		FullName:  p.FullName,
		Phone:     p.Phone,
		Apartment: p.Apartment,
		HouseID:   p.HouseID,
		UKID:      p.UKID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func companyToDomain(c *dbManagementCompany) *domain.ManagementCompany {
	return &domain.ManagementCompany{ID: c.ID, Name: c.Name}
}

func houseToDomain(h *dbHouse) *domain.House {
	return &domain.House{ID: h.ID, Name: h.Name, Address: h.Address, UKID: h.UKID}
}
