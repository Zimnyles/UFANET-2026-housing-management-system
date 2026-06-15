package resources

import (
	"api-gateway/infra/clients/auth_client"
	"api-gateway/pkg/logger"
	"fmt"

	"github.com/gofiber/fiber/v2"
	redisstore "github.com/gofiber/storage/redis/v2"
	"github.com/rs/zerolog"
)

type Resources struct {
	Env        *Env
	Logger     *zerolog.Logger
	Cache      fiber.Storage
	AuthClient *auth_client.AuthClient
}

func InitResources() (*Resources, error) {
	env, err := initEnv()
	if err != nil {
		return nil, err
	}

	log := logger.NewLogger(env.ServiceName, env.LogLevel)
	log.Info().Str("addr", env.Addr()).Msg("env loaded")

	cache, err := initCache(env)
	if err != nil {
		return nil, err
	}
	log.Info().Str("host", env.RedisHost).Int("port", env.RedisPort).Msg("cache connected")

	authClient, err := auth_client.New(env.AuthAddr, log)
	if err != nil {
		return nil, err
	}

	return &Resources{
		Env:        env,
		Logger:     log,
		Cache:      cache,
		AuthClient: authClient,
	}, nil
}

func initCache(env *Env) (*redisstore.Storage, error) {
	store := redisstore.New(redisstore.Config{
		Host:     env.RedisHost,
		Port:     env.RedisPort,
		Password: env.RedisPassword,
		Database: env.RedisDB,
	})

	if err := store.Set("ping", []byte("pong"), 0); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	_ = store.Delete("ping")

	return store, nil
}
