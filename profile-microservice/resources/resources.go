package resources

import (
	"context"
	"fmt"
	repo "profile-service/infra/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Resources struct {
	Env    *Env
	Logger *zerolog.Logger
	Repo   *repo.Repository
}

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
