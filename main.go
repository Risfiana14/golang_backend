package main

import (
	"log"
	"tugas5/config"
	"tugas5/database"
	"tugas5/routes"

	"github.com/gofiber/fiber/v2"
	
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	database.ConnectDB()
	defer database.DB.Close()

	// Fiber app dengan custom error handler
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware global
	app.Use(config.LoggerMiddleware())

	// Setup routes
	routes.UserRoutes(app)

	// Get port dari env (default 3000)
	port := config.GetEnv("APP_PORT", "3000")

	log.Println("Server running on port " + port)
	log.Fatal(app.Listen(":" + port))

}

