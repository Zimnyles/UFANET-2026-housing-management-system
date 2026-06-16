package news_handler

import (
	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type NewsHandler struct {
	logger *zerolog.Logger
}

func NewHandler(logger *zerolog.Logger) *NewsHandler {
	return &NewsHandler{logger: logger}
}

func (h *NewsHandler) List(c *fiber.Ctx) error {
	return app_errors.Respond(c, app_errors.ErrNotImplemented)
}

func (h *NewsHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateNewsRequest
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	if err := validateCreateNewsRequest(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	return app_errors.Respond(c, app_errors.ErrNotImplemented)
}
