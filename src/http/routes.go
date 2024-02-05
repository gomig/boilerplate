package http

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes register web routes
func RegisterRoutes(router *fiber.App) {
	// Routes
	router.Get("/thumb/*", Thumbnail)
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Mig app")
	})

	// Fallback
	router.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})
}
