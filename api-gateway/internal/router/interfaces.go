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

type ProfileHandler interface {
	GetProfile(c *fiber.Ctx) error
	UpsertProfile(c *fiber.Ctx) error
	ListManagementCompanies(c *fiber.Ctx) error
	CreateManagementCompany(c *fiber.Ctx) error
	ListHouses(c *fiber.Ctx) error
	CreateHouse(c *fiber.Ctx) error
}

type RequestsHandler interface {
	Create(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	UpdateStatus(c *fiber.Ctx) error
	AddComment(c *fiber.Ctx) error
	GetComments(c *fiber.Ctx) error
}
