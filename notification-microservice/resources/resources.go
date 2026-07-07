package resources

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"notification-service/infra/repository"
)

type Resources struct {
	Env    *Env
	Logger *zerolog.Logger
	Repo   *repository.Repository
	DB     *gorm.DB
}

func InitResources(ctx context.Context) (*Resources, error) {
	env, err := initEnv()
	if err != nil {
		return nil, err
	}

	level, err := zerolog.ParseLevel(env.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	logger := log.Level(level).With().Str("service", env.ServiceName).Logger()

	db, err := gorm.Open(postgres.Open(env.DSN()), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	repo := repository.New(db)
	if err := repo.Migrate(ctx); err != nil {
		return nil, err
	}

	return &Resources{Env: env, Logger: &logger, Repo: repo, DB: db}, nil
}

func (r *Resources) Close() {
	if sqlDB, err := r.DB.DB(); err == nil {
		_ = sqlDB.Close()
	}
}
