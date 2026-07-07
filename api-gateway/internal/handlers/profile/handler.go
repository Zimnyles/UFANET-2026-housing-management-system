package profile_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/constants"
	"api-gateway/internal/models/dto"
)

type ProfileHandler struct {
	service ProfileService
	logger  *zerolog.Logger
}

func NewHandler(service ProfileService, logger *zerolog.Logger) *ProfileHandler {
	return &ProfileHandler{service: service, logger: logger}
}

func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals(constants.LocalUserID).(string)
	if !ok || userID == "" {
		return app_errors.Respond(c, app_errors.ErrUnauthorized)
	}

	p, err := h.service.GetProfile(c.UserContext(), userID)
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.JSON(toDTOProfileResponse(p))
}

func (h *ProfileHandler) UpsertProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals(constants.LocalUserID).(string)
	if !ok || userID == "" {
		return app_errors.Respond(c, app_errors.ErrUnauthorized)
	}

	var req dto.UpsertProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	if err := validateUpsertProfile(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	p, err := h.service.UpsertProfile(c.UserContext(), toDomainUpsertRequest(userID, req))
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.JSON(toDTOProfileResponse(p))
}

func (h *ProfileHandler) ListManagementCompanies(c *fiber.Ctx) error {
	companies, err := h.service.ListManagementCompanies(c.UserContext())
	if err != nil {
		return app_errors.Respond(c, err)
	}

	items := make([]dto.ManagementCompanyResponse, 0, len(companies))
	for _, mc := range companies {
		items = append(items, toDTOCompanyResponse(mc))
	}

	return c.JSON(dto.ListManagementCompaniesResponse{Companies: items})
}

func (h *ProfileHandler) CreateManagementCompany(c *fiber.Ctx) error {
	var req dto.CreateManagementCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	if err := validateCreateCompany(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	mc, err := h.service.CreateManagementCompany(c.UserContext(), toDomainCreateCompanyRequest(req))
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(toDTOCompanyResponse(mc))
}

func (h *ProfileHandler) ListHouses(c *fiber.Ctx) error {
	ukID := c.Query("uk_id")

	houses, err := h.service.ListHouses(c.UserContext(), ukID)
	if err != nil {
		return app_errors.Respond(c, err)
	}

	items := make([]dto.HouseResponse, 0, len(houses))
	for _, house := range houses {
		items = append(items, toDTOHouseResponse(house))
	}

	return c.JSON(dto.ListHousesResponse{Houses: items})
}

func (h *ProfileHandler) CreateHouse(c *fiber.Ctx) error {
	var req dto.CreateHouseRequest
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	if err := validateCreateHouse(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	house, err := h.service.CreateHouse(c.UserContext(), toDomainCreateHouseRequest(req))
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(toDTOHouseResponse(house))
}
