package config

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Setup middleware
	app.Use(LoggerMiddleware())

	return app
}

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("%s %s", c.Method(), c.Path())
		return c.Next()
	}
}
