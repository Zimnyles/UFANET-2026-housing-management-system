package app_errors

import "github.com/gofiber/fiber/v2"

func Respond(c *fiber.Ctx, err error) error {
	return c.Status(StatusCode(err)).JSON(fiber.Map{"error": err.Error()})
}
