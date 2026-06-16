package api

import (
	"context"

	auth_handler "api-gateway/internal/handlers/auth"
	health_handler "api-gateway/internal/handlers/health"
	news_handler "api-gateway/internal/handlers/news"
	notifications_handler "api-gateway/internal/handlers/notifications"
	profile_handler "api-gateway/internal/handlers/profile"
	"api-gateway/internal/middlewares"
	"api-gateway/internal/router"
	"api-gateway/resources"
	auth_service "api-gateway/services/auth"
	profile_service "api-gateway/services/profile"
	"github.com/gofiber/fiber/v2"
)

type API struct {
	res *resources.Resources
}

func New(res *resources.Resources) *API {
	return &API{res: res}
}

func (a *API) Start(ctx context.Context) error {
	app := fiber.New(fiber.Config{
		AppName:       a.res.Env.ServiceName,
		Prefork:       a.res.Env.Prefork,
		CaseSensitive: true,
	})

	mw := middlewares.New(ctx, a.res, app)
	mw.SetGlobalMiddlewares()

	authService := auth_service.New(a.res.AuthClient, a.res.Logger)
	profileService := profile_service.New(a.res.ProfileClient, a.res.Logger)

	handlers := router.Handlers{
		Health:        health_handler.NewHandler(a.res.Env.ServiceName),
		Auth:          auth_handler.NewHandler(authService, mw, a.res.Logger),
		News:          news_handler.NewHandler(a.res.Logger),
		Notifications: notifications_handler.NewHandler(a.res.Logger),
		Profile:       profile_handler.NewHandler(profileService, a.res.Logger),
	}

	router.New(app, mw, handlers).Register()

	errCh := make(chan error, 1)

	go func() {
		a.res.Logger.Info().Str("addr", a.res.Env.Addr()).Msg("starting server")
		errCh <- app.Listen(a.res.Env.Addr())
	}()

	select {
	case <-ctx.Done():
		a.res.Logger.Info().Msg("shutting down server")

		err := app.Shutdown()

		a.res.Close()

		return err
	case err := <-errCh:
		return err
	}
}
