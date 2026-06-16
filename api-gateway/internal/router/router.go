package router

import (
	"api-gateway/internal/middlewares"
	"api-gateway/internal/models/constants"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Auth          AuthHandler
	Health        HealthHandler
	News          NewsHandler
	Notifications NotificationsHandler
	Profile       ProfileHandler
}

type Router struct {
	app *fiber.App
	mw  *middlewares.Middlewares
	h   Handlers
}

func New(app *fiber.App, mw *middlewares.Middlewares, h Handlers) *Router {
	return &Router{app: app, mw: mw, h: h}
}

func (r *Router) Register() {
	r.health()
	r.auth()
	r.news()
	r.notifications()
	r.profile()
}

func (r *Router) health() {
	r.app.Get("/health", r.h.Health.Check)
}

func (r *Router) auth() {
	g := r.app.Group("/auth")
	g.Post("/register", r.mw.RateLimit(), r.h.Auth.Register)
	g.Post("/login", r.mw.RateLimit(), r.h.Auth.Login)
	g.Post("/refresh", r.h.Auth.Refresh)
	g.Post("/logout", r.mw.JWTAuth(), r.h.Auth.Logout)
}

func (r *Router) news() {
	g := r.app.Group("/news")
	g.Use(r.mw.JWTAuth())
	g.Use(r.mw.Timeout())
	g.Get("/", r.h.News.List)
	g.Post("/", r.mw.RequireRole(constants.RoleAdmin), r.h.News.Create)
}

func (r *Router) notifications() {
	g := r.app.Group("/notifications")
	g.Use(r.mw.JWTAuth())
	g.Post("/register", r.h.Notifications.Register)
	g.Delete("/unregister", r.h.Notifications.Unregister)
}

func (r *Router) profile() {
	g := r.app.Group("/profile")
	g.Use(r.mw.JWTAuth())
	g.Use(r.mw.Timeout())
	g.Get("/", r.h.Profile.GetProfile)
	g.Put("/", r.h.Profile.UpsertProfile)

	mc := r.app.Group("/management-companies")
	mc.Use(r.mw.JWTAuth())
	mc.Use(r.mw.Timeout())
	mc.Get("/", r.h.Profile.ListManagementCompanies)
	mc.Post("/", r.mw.RequireRole(constants.RoleAdmin), r.h.Profile.CreateManagementCompany)

	houses := r.app.Group("/houses")
	houses.Use(r.mw.JWTAuth())
	houses.Use(r.mw.Timeout())
	houses.Get("/", r.h.Profile.ListHouses)
	houses.Post("/", r.mw.RequireRole(constants.RoleAdmin), r.h.Profile.CreateHouse)
}
