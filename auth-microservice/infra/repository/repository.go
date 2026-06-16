package repository

import (
	infra_errors "auth-service/infra/errors"
	"auth-service/infra/models/domain"
	"context"
	"errors"
	"fmt"

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
		CREATE TABLE IF NOT EXISTS users (
			id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name       TEXT NOT NULL,
			email      TEXT NOT NULL UNIQUE,
			password   TEXT NOT NULL,
			role       TEXT NOT NULL DEFAULT 'user',
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	var created domain.User
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (name, email, password, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, email, password, role, created_at`,
		user.Name, user.Email, user.PasswordHash, user.Role,
	).Scan(&created.ID, &created.Name, &created.Email, &created.PasswordHash, &created.Role, &created.CreatedAt)
	if err != nil {
		if containsCode(err, "23505") {
			return nil, infra_errors.ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &created, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, email, password, role, created_at
		 FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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
