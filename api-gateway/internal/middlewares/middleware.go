package middlewares

import (
	"api-gateway/resources"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
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

func (mw *Middlewares) SetGlobalMiddlewares() {
	mw.app.Use(recover.New())
	mw.app.Use(mw.requestLogger(mw.logger))
}
