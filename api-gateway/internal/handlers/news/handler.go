package news_handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/constants"
	"api-gateway/internal/models/dto"
)

type NewsHandler struct {
	service  NewsService
	profiles ProfileService
	logger   *zerolog.Logger
}

func NewHandler(service NewsService, profiles ProfileService, logger *zerolog.Logger) *NewsHandler {
	return &NewsHandler{service: service, profiles: profiles, logger: logger}
}

func (h *NewsHandler) List(c *fiber.Ctx) error {
	houseID := ""

	if currentRole(c) != constants.RoleAdmin {
		var err error

		houseID, err = h.currentHouseID(c)
		if err != nil {
			return app_errors.Respond(c, err)
		}
	}

	items, total, err := h.service.List(c.UserContext(), toDomainList(houseID, c.Query("date_from"), c.Query("date_to"), newsQueryInt(c, "limit"), newsQueryInt(c, "offset")))
	if err != nil {
		return app_errors.Respond(c, err)
	}

	responses := make([]dto.NewsResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toDTO(item, h.authorName(c.UserContext(), item.CreatedBy)))
	}

	return c.JSON(dto.ListNewsResponse{News: responses, Total: total})
}

func (h *NewsHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateNewsRequest
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	if err := validateCreateNewsRequest(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	userID, ok := c.Locals(constants.LocalUserID).(string)
	if !ok || userID == "" {
		return app_errors.Respond(c, app_errors.ErrUnauthorized)
	}

	houseID := req.HouseID
	if houseID == "" {
		var err error

		houseID, err = h.currentHouseID(c)
		if err != nil {
			return app_errors.Respond(c, err)
		}
	}

	if !h.houseExists(c.UserContext(), houseID) {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	news, err := h.service.Create(c.UserContext(), toDomainCreate(req, houseID, userID))
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(toDTO(news, h.authorName(c.UserContext(), news.CreatedBy)))
}

func (h *NewsHandler) houseExists(ctx context.Context, houseID string) bool {
	houses, err := h.profiles.ListHouses(ctx, "")
	if err != nil {
		return false
	}

	for _, house := range houses {
		if house.ID == houseID {
			return true
		}
	}

	return false
}

func (h *NewsHandler) authorName(ctx context.Context, userID string) string {
	profile, err := h.profiles.GetProfile(ctx, userID)
	if err != nil || profile == nil || profile.FullName == "" {
		return "Пользователь"
	}

	return profile.FullName
}

func currentRole(c *fiber.Ctx) string {
	role, _ := c.Locals(constants.LocalRole).(string)

	return role
}

func (h *NewsHandler) currentHouseID(c *fiber.Ctx) (string, error) {
	userID, ok := c.Locals(constants.LocalUserID).(string)
	if !ok || userID == "" {
		return "", app_errors.ErrUnauthorized
	}

	profile, err := h.profiles.GetProfile(c.UserContext(), userID)
	if err != nil {
		return "", err
	}

	if profile == nil || profile.HouseID == "" {
		return "", app_errors.ErrBadRequest
	}

	return profile.HouseID, nil
}

func newsQueryInt(c *fiber.Ctx, key string) int {
	value, err := strconv.Atoi(c.Query(key))
	if err != nil {
		return 0
	}

	return value
}
