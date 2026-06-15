package auth_handler

import (
	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	authService AuthService
	logger      *zerolog.Logger
}

func NewHandler(authService AuthService, logger *zerolog.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, logger: logger}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}
	if err := validateRegisterRequest(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	result, err := h.authService.Register(c.Context(), toDomainRegisterRequest(req))
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

	result, err := h.authService.Login(c.Context(), toDomainLoginRequest(req))
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

	tokens, err := h.authService.Refresh(c.Context(), toDomainRefreshRequest(req))
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

	if err := h.authService.Logout(c.Context(), req.RefreshToken); err != nil {
		return app_errors.Respond(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
