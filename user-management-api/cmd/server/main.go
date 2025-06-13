package main

import (
	"log"
	"user-management-api/internal/database"
	"user-management-api/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Initialize database
	database.InitDatabase()

	// Create Express.js server with custom configuration
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Bug: Generic error handling - could leak internal details
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Public routes
	app.Post("/register", handlers.RegisterUser)
	app.Post("/login", handlers.LoginUser)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Protected routes group
	api := app.Group("/api/v1", handlers.AuthMiddleware)
	api.Post("/updateUser/:id", handlers.UpdateUser)

	// Admin only routes
	admin := api.Group("/admin", handlers.AdminMiddleware)
	admin.Get("/users", handlers.GetUsers)

	log.Println("Server starting on :8080...")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
