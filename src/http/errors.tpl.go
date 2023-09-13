package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gomig/http"
)

// OnError handle web app errors and generate good response
func OnError(c *fiber.Ctx, err error) error {
	if http.WantJson(c) {
		return c.SendStatus(500)
	} else {
		return c.Status(500).SendString("Error On server")
	}
}
