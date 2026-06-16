package repository

import (
	"context"
	"errors"
	"fmt"

	infra_errors "profile-service/infra/errors"
	"profile-service/infra/models/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Migrate(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS management_companies (
			id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name       TEXT NOT NULL UNIQUE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);

		CREATE TABLE IF NOT EXISTS houses (
			id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name       TEXT NOT NULL,
			address    TEXT NOT NULL,
			uk_id      UUID NOT NULL REFERENCES management_companies(id),
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);

		CREATE TABLE IF NOT EXISTS profiles (
			user_id    UUID PRIMARY KEY,
			full_name  TEXT NOT NULL DEFAULT '',
			phone      TEXT NOT NULL DEFAULT '',
			apartment  TEXT NOT NULL DEFAULT '',
			house_id   UUID REFERENCES houses(id),
			uk_id      UUID REFERENCES management_companies(id),
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

// ─── profile ──────────────────────────────────────────────────────────────────

func (r *Repository) GetProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	var p domain.Profile
	err := r.pool.QueryRow(ctx,
		`SELECT user_id, full_name, phone, apartment,
		        COALESCE(house_id::text, ''), COALESCE(uk_id::text, ''),
		        created_at, updated_at
		 FROM profiles WHERE user_id = $1`,
		userID,
	).Scan(&p.UserID, &p.FullName, &p.Phone, &p.Apartment, &p.HouseID, &p.UKID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, infra_errors.ErrProfileNotFound
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return &p, nil
}

func (r *Repository) UpsertProfile(ctx context.Context, profile *domain.Profile) (*domain.Profile, error) {
	var p domain.Profile
	var houseID *string
	if profile.HouseID != "" {
		houseID = &profile.HouseID
	}

	// uk_id подтягивается из house
	err := r.pool.QueryRow(ctx,
		`INSERT INTO profiles (user_id, full_name, phone, apartment, house_id, uk_id)
		 VALUES ($1, $2, $3, $4, $5,
		   (SELECT uk_id FROM houses WHERE id = $5))
		 ON CONFLICT (user_id) DO UPDATE SET
		   full_name  = EXCLUDED.full_name,
		   phone      = EXCLUDED.phone,
		   apartment  = EXCLUDED.apartment,
		   house_id   = EXCLUDED.house_id,
		   uk_id      = EXCLUDED.uk_id,
		   updated_at = now()
		 RETURNING user_id, full_name, phone, apartment,
		           COALESCE(house_id::text, ''), COALESCE(uk_id::text, ''),
		           created_at, updated_at`,
		profile.UserID, profile.FullName, profile.Phone, profile.Apartment, houseID,
	).Scan(&p.UserID, &p.FullName, &p.Phone, &p.Apartment, &p.HouseID, &p.UKID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("upsert profile: %w", err)
	}
	return &p, nil
}

func (r *Repository) IsProfileComplete(ctx context.Context, userID string) (bool, error) {
	var complete bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM profiles
			WHERE user_id = $1
			  AND full_name <> ''
			  AND phone     <> ''
			  AND apartment <> ''
			  AND house_id IS NOT NULL
		)`,
		userID,
	).Scan(&complete)
	if err != nil {
		return false, fmt.Errorf("is profile complete: %w", err)
	}
	return complete, nil
}

// ─── management companies ─────────────────────────────────────────────────────

func (r *Repository) CreateManagementCompany(ctx context.Context, name string) (*domain.ManagementCompany, error) {
	var c domain.ManagementCompany
	err := r.pool.QueryRow(ctx,
		`INSERT INTO management_companies (name)
		 VALUES ($1)
		 RETURNING id, name`,
		name,
	).Scan(&c.ID, &c.Name)
	if err != nil {
		return nil, fmt.Errorf("create management company: %w", err)
	}
	return &c, nil
}

func (r *Repository) ListManagementCompanies(ctx context.Context) ([]*domain.ManagementCompany, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name FROM management_companies ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("list management companies: %w", err)
	}
	defer rows.Close()

	var companies []*domain.ManagementCompany
	for rows.Next() {
		var c domain.ManagementCompany
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, fmt.Errorf("scan management company: %w", err)
		}
		companies = append(companies, &c)
	}
	return companies, nil
}

// ─── houses ───────────────────────────────────────────────────────────────────

func (r *Repository) CreateHouse(ctx context.Context, house *domain.House) (*domain.House, error) {
	var h domain.House
	err := r.pool.QueryRow(ctx,
		`INSERT INTO houses (name, address, uk_id)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, address, uk_id::text`,
		house.Name, house.Address, house.UKID,
	).Scan(&h.ID, &h.Name, &h.Address, &h.UKID)
	if err != nil {
		return nil, fmt.Errorf("create house: %w", err)
	}
	return &h, nil
}

func (r *Repository) ListHouses(ctx context.Context, ukID string) ([]*domain.House, error) {
	var (
		rows pgx.Rows
		err  error
	)
	if ukID != "" {
		rows, err = r.pool.Query(ctx,
			`SELECT id, name, address, uk_id::text FROM houses WHERE uk_id = $1 ORDER BY name`,
			ukID,
		)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT id, name, address, uk_id::text FROM houses ORDER BY name`,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("list houses: %w", err)
	}
	defer rows.Close()

	var houses []*domain.House
	for rows.Next() {
		var h domain.House
		if err := rows.Scan(&h.ID, &h.Name, &h.Address, &h.UKID); err != nil {
			return nil, fmt.Errorf("scan house: %w", err)
		}
		houses = append(houses, &h)
	}
	return houses, nil
}
