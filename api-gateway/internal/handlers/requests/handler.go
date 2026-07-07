package requests_handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/constants"
	"api-gateway/internal/models/domain"
	"api-gateway/internal/models/dto"
)

type RequestsHandler struct {
	service  RequestsService
	profiles ProfileService
	logger   *zerolog.Logger
}

func NewHandler(service RequestsService, profiles ProfileService, logger *zerolog.Logger) *RequestsHandler {
	return &RequestsHandler{service: service, profiles: profiles, logger: logger}
}

func (h *RequestsHandler) Create(c *fiber.Ctx) error {
	userID, ok := currentUserID(c)
	if !ok {
		return app_errors.Respond(c, app_errors.ErrUnauthorized)
	}

	var req dto.CreateMaintenanceRequest
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	if err := validateCreateRequest(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	r, err := h.service.CreateRequest(c.UserContext(), toDomainCreateRequest(userID, req))
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(h.toDTORequest(c.UserContext(), r))
}

func (h *RequestsHandler) List(c *fiber.Ctx) error {
	userID, ok := currentUserID(c)
	if !ok {
		return app_errors.Respond(c, app_errors.ErrUnauthorized)
	}

	filterUserID := userID
	if currentRole(c) == constants.RoleAdmin {
		filterUserID = c.Query("user_id")
	}

	req := toDomainListRequest(filterUserID, c.Query("status"), c.Query("type"), queryInt(c, "limit"), queryInt(c, "offset"))

	items, total, err := h.service.GetRequests(c.UserContext(), req)
	if err != nil {
		return app_errors.Respond(c, err)
	}

	responses := make([]dto.MaintenanceRequestResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, h.toDTORequest(c.UserContext(), item))
	}

	return c.JSON(dto.ListMaintenanceRequestsResponse{Requests: responses, Total: total})
}

func (h *RequestsHandler) Get(c *fiber.Ctx) error {
	r, err := h.getVisibleRequest(c)
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.JSON(h.toDTORequest(c.UserContext(), r))
}

func (h *RequestsHandler) UpdateStatus(c *fiber.Ctx) error {
	userID, ok := currentUserID(c)
	if !ok {
		return app_errors.Respond(c, app_errors.ErrUnauthorized)
	}

	id := c.Params("id")
	if id == "" {
		return app_errors.Respond(c, app_errors.ErrRequestIDRequired)
	}

	var req dto.UpdateMaintenanceRequestStatus
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	if err := validateUpdateStatus(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	r, err := h.service.UpdateRequestStatus(c.UserContext(), toDomainUpdateStatus(id, userID, req))
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.JSON(h.toDTORequest(c.UserContext(), r))
}

func (h *RequestsHandler) AddComment(c *fiber.Ctx) error {
	r, err := h.getVisibleRequest(c)
	if err != nil {
		return app_errors.Respond(c, err)
	}

	userID, _ := currentUserID(c)

	var req dto.AddMaintenanceRequestComment
	if err := c.BodyParser(&req); err != nil {
		return app_errors.Respond(c, app_errors.ErrBadRequest)
	}

	if err := validateAddComment(&req); err != nil {
		return app_errors.Respond(c, err)
	}

	comment, err := h.service.AddComment(c.UserContext(), toDomainAddComment(r.ID, userID, req))
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(h.toDTOComment(c.UserContext(), comment))
}

func (h *RequestsHandler) GetComments(c *fiber.Ctx) error {
	r, err := h.getVisibleRequest(c)
	if err != nil {
		return app_errors.Respond(c, err)
	}

	items, err := h.service.GetComments(c.UserContext(), r.ID)
	if err != nil {
		return app_errors.Respond(c, err)
	}

	responses := make([]dto.RequestCommentResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, h.toDTOComment(c.UserContext(), item))
	}

	return c.JSON(dto.ListMaintenanceRequestCommentsResponse{Comments: responses})
}

func (h *RequestsHandler) toDTORequest(ctx context.Context, item *domain.MaintenanceRequest) dto.MaintenanceRequestResponse {
	response := toDTORequest(item)
	response.AuthorName = h.authorName(ctx, item.UserID)

	return response
}

func (h *RequestsHandler) toDTOComment(ctx context.Context, item *domain.RequestComment) dto.RequestCommentResponse {
	response := toDTOComment(item)
	response.AuthorName = h.authorName(ctx, item.UserID)

	return response
}

func (h *RequestsHandler) authorName(ctx context.Context, userID string) string {
	profile, err := h.profiles.GetProfile(ctx, userID)
	if err != nil || profile == nil || profile.FullName == "" {
		return "Пользователь"
	}

	return profile.FullName
}

func (h *RequestsHandler) getVisibleRequest(c *fiber.Ctx) (*domain.MaintenanceRequest, error) {
	userID, ok := currentUserID(c)
	if !ok {
		return nil, app_errors.ErrUnauthorized
	}

	id := c.Params("id")
	if id == "" {
		return nil, app_errors.ErrRequestIDRequired
	}

	r, err := h.service.GetRequest(c.UserContext(), id)
	if err != nil {
		return nil, err
	}

	if currentRole(c) != constants.RoleAdmin && r.UserID != userID {
		return nil, app_errors.ErrForbidden
	}

	return r, nil
}

func currentUserID(c *fiber.Ctx) (string, bool) {
	userID, ok := c.Locals(constants.LocalUserID).(string)

	return userID, ok && userID != ""
}

func currentRole(c *fiber.Ctx) string {
	role, _ := c.Locals(constants.LocalRole).(string)

	return role
}

func queryInt(c *fiber.Ctx, key string) int {
	value := c.Query(key)
	if value == "" {
		return 0
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return n
}
