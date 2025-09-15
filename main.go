package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"project/database"
	"project/routes"
)
func main() {
	godotenv.Load()

	database.ConnectDb()

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Go-Mart!")
	})

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}