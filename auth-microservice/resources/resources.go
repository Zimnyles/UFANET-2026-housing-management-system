package resources

import (
	"auth-service/infra/hasher"
	"auth-service/infra/jwt"
	"auth-service/infra/models/domain"
	repo "auth-service/infra/repository"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Resources struct {
	Env    *Env
	Logger *zerolog.Logger
	Repo   *repo.Repository
	JWT    *jwtAdapter
	Hasher *hasherAdapter
}

type jwtAdapter struct {
	m *jwt.Manager
}

func (a *jwtAdapter) GenerateAccess(userID, role string) (string, error) {
	return a.m.GenerateAccess(userID, role)
}

func (a *jwtAdapter) GenerateRefresh(userID, role string) (string, error) {
	return a.m.GenerateRefresh(userID, role)
}

func (a *jwtAdapter) ParseRefresh(tokenStr string) (*domain.TokenClaims, error) {
	claims, err := a.m.ParseRefresh(tokenStr)
	if err != nil {
		return nil, err
	}
	return &domain.TokenClaims{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}

type hasherAdapter struct{}

func (h *hasherAdapter) Hash(password string) (string, error) { return hasher.Hash(password) }
func (h *hasherAdapter) Check(password, hash string) bool     { return hasher.Check(password, hash) }

func InitResources(ctx context.Context) (*Resources, error) {
	env, err := initEnv()
	if err != nil {
		return nil, err
	}

	logger := initLogger(env.ServiceName, env.LogLevel)
	logger.Info().Str("addr", env.Addr()).Msg("env loaded")

	pool, err := pgxpool.New(ctx, env.DSN())
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	logger.Info().Str("host", env.PostgresHost).Int("port", env.PostgresPort).Msg("postgres connected")

	repository := repo.New(pool)
	if err := repository.Migrate(ctx); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	logger.Info().Msg("migrations applied")

	return &Resources{
		Env:    env,
		Logger: logger,
		Repo:   repository,
		JWT:    &jwtAdapter{m: jwt.NewManager(env.JWTSecret, env.JWTAccessTTL, env.JWTRefreshTTL)},
		Hasher: &hasherAdapter{},
	}, nil
}

func initLogger(serviceName, level string) *zerolog.Logger {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	logger := log.Level(lvl).With().Str("service", serviceName).Logger()
	return &logger
}
