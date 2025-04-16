package routes

import(
	"github.com/gofiber/fiber/v2"
	"github.com/kithmina1999/eldraread-api/controllers"
)

func RegisterUserRoutes(app *fiber.App){
	api := app.Group("/api/users")
	api.Post("/register",controllers.RegisterUser)
	api.Post("/register/google",controllers.RegisterWithGoogle)
	api.Post("/login",controllers.LoginUser)
}