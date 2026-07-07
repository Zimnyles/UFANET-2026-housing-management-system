package notifications_handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/constants"
	"api-gateway/internal/models/dto"
)

type NotificationsHandler struct {
	service  NotificationsService
	profiles ProfileService
	logger   *zerolog.Logger
}

func NewHandler(service NotificationsService, profiles ProfileService, logger *zerolog.Logger) *NotificationsHandler {
	return &NotificationsHandler{service: service, profiles: profiles, logger: logger}
}

func (h *NotificationsHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	if err := validateRegisterDeviceRequest(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	userID, ok := c.Locals(constants.LocalUserID).(string)
	if !ok || userID == "" {
		return app_errors.Respond(c, app_errors.ErrUnauthorized)
	}

	if err := h.service.Register(c.UserContext(), userID, toDomainRegister(req)); err != nil {
		return app_errors.Respond(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *NotificationsHandler) Unregister(c *fiber.Ctx) error {
	userID, ok := c.Locals(constants.LocalUserID).(string)
	if !ok || userID == "" {
		return app_errors.Respond(c, app_errors.ErrUnauthorized)
	}

	token := c.Query("device_token")
	if token == "" {
		return app_errors.Respond(c, app_errors.ErrDeviceTokenRequired)
	}

	if err := h.service.Unregister(c.UserContext(), userID, token); err != nil {
		return app_errors.Respond(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *NotificationsHandler) Stream(c *fiber.Ctx) error {
	userID, ok := c.Locals(constants.LocalUserID).(string)
	if !ok || userID == "" {
		return app_errors.Respond(c, app_errors.ErrUnauthorized)
	}

	var houseID string

	profile, err := h.profiles.GetProfile(c.UserContext(), userID)
	if err == nil && profile != nil {
		houseID = profile.HouseID
	}

	events, unsubscribe, err := h.service.Subscribe(context.Background(), userID, houseID)
	if err != nil {
		return app_errors.Respond(c, err)
	}

	c.Set(fiber.HeaderContentType, "text/event-stream; charset=utf-8")
	c.Set(fiber.HeaderCacheControl, "no-cache, no-transform")
	c.Set(fiber.HeaderConnection, "keep-alive")
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer unsubscribe()

		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		if _, err := fmt.Fprint(w, ": connected\n\n"); err != nil {
			return
		}

		if err := w.Flush(); err != nil {
			return
		}

		for {
			select {
			case event, open := <-events:
				if !open {
					return
				}

				payload, err := json.Marshal(event)
				if err != nil {
					continue
				}

				if _, err := fmt.Fprintf(w, "event: notification\ndata: %s\n\n", payload); err != nil {
					return
				}
			case <-ticker.C:
				if _, err := fmt.Fprint(w, ": heartbeat\n\n"); err != nil {
					return
				}
			}

			if err := w.Flush(); err != nil {
				return
			}
		}
	})

	return nil
}
