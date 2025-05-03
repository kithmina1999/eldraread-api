package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/kithmina1999/eldraread-api/db"
	"github.com/kithmina1999/eldraread-api/routes"
)

func main() {
	godotenv.Load()
	db.Connect()

	app := fiber.New()


	app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:3000, http://localhost:3001", // Allow both origins
        AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders: "Content-Type,Authorization",
        AllowCredentials: true,
    }))

	routes.RegisterUserRoutes(app)
	// routes.RegisteredUserRoutes(app)
	routes.AdminRoutes(app)

	app.Listen(":8080")
}
