package resources

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"news-service/infra/mq"
	repo "news-service/infra/repository"
)

type Resources struct {
	Env       *Env
	Logger    *zerolog.Logger
	Repo      *repo.Repository
	Publisher *mq.Publisher
}

func InitResources(ctx context.Context) (*Resources, error) {
	env, err := initEnv()
	if err != nil {
		return nil, err
	}

	logger := initLogger(env.ServiceName, env.LogLevel)
	logger.Info().Str("addr", env.Addr()).Msg("env loaded")

	db, err := gorm.Open(postgres.Open(env.DSN()), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("postgres handle: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	logger.Info().Str("host", env.PostgresHost).Int("port", env.PostgresPort).Msg("postgres connected")

	repository := repo.New(db)
	if err := repository.Migrate(ctx); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	logger.Info().Msg("migrations applied")

	publisher, err := mq.New(env.RabbitDSN(), logger)
	if err != nil {
		return nil, fmt.Errorf("init mq: %w", err)
	}

	return &Resources{
		Env:       env,
		Logger:    logger,
		Repo:      repository,
		Publisher: publisher,
	}, nil
}

func (r *Resources) Close() {
	if r.Publisher != nil {
		if err := r.Publisher.Close(); err != nil {
			r.Logger.Error().Err(err).Msg("failed to close publisher")
		}
	}
}

func initLogger(serviceName, level string) *zerolog.Logger {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	logger := log.Level(lvl).With().Str("service", serviceName).Logger()

	return &logger
}
