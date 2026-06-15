package notifications_handler

import (
	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type NotificationsHandler struct {
	logger *zerolog.Logger
}

func NewHandler(logger *zerolog.Logger) *NotificationsHandler {
	return &NotificationsHandler{logger: logger}
}

func (h *NotificationsHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}
	if err := validateRegisterDeviceRequest(&req); err != nil {
		return app_errors.Respond(c, err)
	}
	return app_errors.Respond(c, app_errors.ErrNotImplemented)
}

func (h *NotificationsHandler) Unregister(c *fiber.Ctx) error {
	return app_errors.Respond(c, app_errors.ErrNotImplemented)
}
