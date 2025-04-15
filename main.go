package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kithmina1999/eldraread-api/routes"
)

func main(){
	app:= fiber.New()

	routes.RegisterUserRoutes(app)

	app.Listen(":8080")
}

