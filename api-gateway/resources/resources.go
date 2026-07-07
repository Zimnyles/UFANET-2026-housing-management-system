package resources

import (
	"github.com/gofiber/fiber/v2"
	redisstore "github.com/gofiber/storage/redis/v2"
	"github.com/rs/zerolog"

	"api-gateway/infra/clients/auth_client"
	"api-gateway/infra/clients/news_client"
	"api-gateway/infra/clients/notifications_client"
	"api-gateway/infra/clients/profile_client"
	"api-gateway/infra/clients/requests_client"
	"api-gateway/pkg/logger"
)

type Resources struct {
	Env                 *Env
	Logger              *zerolog.Logger
	Cache               fiber.Storage
	AuthClient          *auth_client.AuthClient
	ProfileClient       *profile_client.ProfileClient
	RequestsClient      *requests_client.RequestsClient
	NewsClient          *news_client.Client
	NotificationsClient *notifications_client.Client
}

func InitResources() (*Resources, error) {
	env, err := initEnv()
	if err != nil {
		return nil, err
	}

	logger := logger.NewLogger(env.ServiceName, env.LogLevel)
	logger.Info().Str("addr", env.Addr()).Msg("env loaded")

	cache := redisstore.New(redisstore.Config{
		Host:     env.RedisHost,
		Port:     env.RedisPort,
		Password: env.RedisPassword,
		Database: env.RedisDB,
	})

	logger.Info().Str("host", env.RedisHost).Int("port", env.RedisPort).Msg("cache connected")

	authClient, err := auth_client.New(env.AuthAddr, logger)
	if err != nil {
		return nil, err
	}

	profileClient, err := profile_client.New(env.ProfileAddr, logger)
	if err != nil {
		return nil, err
	}

	requestsClient, err := requests_client.New(env.RequestsAddr, logger)
	if err != nil {
		return nil, err
	}

	newsClient, err := news_client.New(env.NewsAddr, logger)
	if err != nil {
		return nil, err
	}

	notificationsClient, err := notifications_client.New(env.NotificationsAddr)
	if err != nil {
		return nil, err
	}

	return &Resources{
		Env:                 env,
		Logger:              logger,
		Cache:               cache,
		AuthClient:          authClient,
		ProfileClient:       profileClient,
		RequestsClient:      requestsClient,
		NewsClient:          newsClient,
		NotificationsClient: notificationsClient,
	}, nil
}

func (r *Resources) Close() {
	if err := r.AuthClient.Close(); err != nil {
		r.Logger.Error().Err(err).Msg("failed to close auth client")
	}

	if err := r.ProfileClient.Close(); err != nil {
		r.Logger.Error().Err(err).Msg("failed to close profile client")
	}

	if err := r.RequestsClient.Close(); err != nil {
		r.Logger.Error().Err(err).Msg("failed to close requests client")
	}

	if err := r.NewsClient.Close(); err != nil {
		r.Logger.Error().Err(err).Msg("failed to close news client")
	}

	if err := r.NotificationsClient.Close(); err != nil {
		r.Logger.Error().Err(err).Msg("failed to close notifications client")
	}

	if err := r.Cache.Reset(); err != nil {
		r.Logger.Error().Err(err).Msg("failed to close cache")
	}
}
