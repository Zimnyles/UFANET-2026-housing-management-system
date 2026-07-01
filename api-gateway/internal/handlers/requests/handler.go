package requests_handler

import (
	"strconv"

	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/constants"
	"api-gateway/internal/models/domain"
	"api-gateway/internal/models/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type RequestsHandler struct {
	service RequestsService
	logger  *zerolog.Logger
}

func NewHandler(service RequestsService, logger *zerolog.Logger) *RequestsHandler {
	return &RequestsHandler{service: service, logger: logger}
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

	return c.Status(fiber.StatusCreated).JSON(toDTORequest(r))
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

	req := &domain.ListMaintenanceRequests{
		UserID: filterUserID,
		Status: c.Query("status"),
		Type:   c.Query("type"),
		Limit:  queryInt(c, "limit"),
		Offset: queryInt(c, "offset"),
	}

	items, total, err := h.service.GetRequests(c.UserContext(), req)
	if err != nil {
		return app_errors.Respond(c, err)
	}

	resp := make([]dto.MaintenanceRequestResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, toDTORequest(item))
	}

	return c.JSON(dto.ListMaintenanceRequestsResponse{Requests: resp, Total: total})
}

func (h *RequestsHandler) Get(c *fiber.Ctx) error {
	r, err := h.getVisibleRequest(c)
	if err != nil {
		return app_errors.Respond(c, err)
	}
	return c.JSON(toDTORequest(r))
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

	r, err := h.service.UpdateRequestStatus(c.UserContext(), &domain.UpdateMaintenanceRequestStatus{
		ID:     id,
		Status: req.Status,
		UserID: userID,
	})
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.JSON(toDTORequest(r))
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

	comment, err := h.service.AddComment(c.UserContext(), &domain.AddMaintenanceRequestComment{
		RequestID: r.ID,
		UserID:    userID,
		Content:   req.Content,
	})
	if err != nil {
		return app_errors.Respond(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(toDTOComment(comment))
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

	resp := make([]dto.RequestCommentResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, toDTOComment(item))
	}

	return c.JSON(dto.ListMaintenanceRequestCommentsResponse{Comments: resp})
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
