package http

import (
	"__ns__/src/app"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gomig/http/middlewares"
	"github.com/gomig/utils"
)

// RegisterGlobalMiddlewares register global app middlewares routes
func RegisterGlobalMiddlewares(router *fiber.App) {
	// Global Middlewares
	router.Use(func(c *fiber.Ctx) error {
		if ok := utils.VarOrPanic(app.IsUnderMaintenance()); ok {
			return c.SendStatus(fiber.StatusServiceUnavailable)
		}
		return c.Next()
	}) // Maintenance mode
	router.Use(middlewares.RateLimiter(
		"GLOBAL-LIMITER",
		app.Config().Cast("web.limit").UInt32Safe(60),
		time.Minute,
		app.Cache(), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		[]string{}, // methods to include only
		[]string{}, // path patterns to ignore
	)) // Accept 60 request in minutes
}
