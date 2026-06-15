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
