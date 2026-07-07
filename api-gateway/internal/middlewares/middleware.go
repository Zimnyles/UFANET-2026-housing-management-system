package middlewares

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"

	"api-gateway/resources"
)

type Middlewares struct {
	jwtSecret           string
	rateLimitMax        int
	rateLimitExpiration time.Duration
	requestTimeout      time.Duration
	storage             fiber.Storage
	logger              *zerolog.Logger
	app                 *fiber.App
}

func New(res *resources.Resources, app *fiber.App) *Middlewares {
	return &Middlewares{
		jwtSecret:           res.Env.JWTSecret,
		rateLimitMax:        res.Env.RateLimitMax,
		rateLimitExpiration: res.Env.RateLimitExpiration,
		requestTimeout:      res.Env.RequestTimeout,
		storage:             res.Cache,
		logger:              res.Logger,
		app:                 app,
	}
}

func (mw *Middlewares) SetGlobalMiddlewares(ctx context.Context) {
	mw.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	mw.app.Use(func(c *fiber.Ctx) error {
		reqCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		c.SetUserContext(reqCtx)

		return c.Next()
	})
	mw.app.Use(recover.New())
	mw.app.Use(mw.requestLogger(mw.logger))
}
