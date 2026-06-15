package health_handler

import "github.com/gofiber/fiber/v2"

type HealthHandler struct {
	serviceName string
}

func NewHandler(serviceName string) *HealthHandler {
	return &HealthHandler{serviceName: serviceName}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": h.serviceName,
	})
}
