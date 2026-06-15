package router

import "github.com/gofiber/fiber/v2"

type AuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type HealthHandler interface {
	Check(c *fiber.Ctx) error
}

type NewsHandler interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
}

type NotificationsHandler interface {
	Register(c *fiber.Ctx) error
	Unregister(c *fiber.Ctx) error
}
