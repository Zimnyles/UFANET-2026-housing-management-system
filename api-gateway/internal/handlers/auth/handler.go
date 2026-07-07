package auth_handler

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/constants"
	"api-gateway/internal/models/dto"
)

type AuthHandler struct {
	authService AuthService
	blacklist   Blacklist
	logger      *zerolog.Logger
}

func NewHandler(authService AuthService, blacklist Blacklist, logger *zerolog.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, blacklist: blacklist, logger: logger}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	if err := validateRegisterRequest(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	result, err := h.authService.Register(c.UserContext(), toDomainRegisterRequest(req))
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(toDTORegisterResponse(result))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	if err := validateLoginRequest(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	result, err := h.authService.Login(c.UserContext(), toDomainLoginRequest(req))
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.JSON(toDTOLoginResponse(result))
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	if err := validateRefreshRequest(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	if h.blacklist.IsBlacklisted(req.RefreshToken) {
		return app_errors.Respond(c, app_errors.ErrInvalidToken)
	}

	tokens, err := h.authService.Refresh(c.UserContext(), toDomainRefreshRequest(req))
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.JSON(toDTORefreshResponse(tokens))
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	if err := validateRefreshRequest(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	if header := c.Get("Authorization"); header != "" {
		if parts := strings.SplitN(header, " ", 2); len(parts) == 2 {
			expiry, ok := c.Locals(constants.LocalTokenExpiry).(time.Time)
			if ok {
				if err := h.blacklist.BlacklistToken(parts[1], time.Until(expiry)); err != nil {
					h.logger.Warn().Err(err).Msg("logout: failed to blacklist access token")
				}
			}
		}
	}

	if err := h.blacklist.BlacklistRawJWT(req.RefreshToken); err != nil {
		h.logger.Warn().Err(err).Msg("logout: failed to blacklist refresh token")
	}

	if err := h.authService.Logout(c.UserContext(), req.RefreshToken); err != nil {
		return app_errors.Respond(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
