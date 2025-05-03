package routes

import(
	"github.com/gofiber/fiber/v2"
	"github.com/kithmina1999/eldraread-api/controllers"
	"github.com/kithmina1999/eldraread-api/middlewares"
)

func AdminRoutes(app *fiber.App){
	api := app.Group("/api/admin")
	api.Post("/login",controllers.AdminLogin)

	//check admin session
	app.Get("/api/session-check",middlewares.AdminSessionCheck)

	//protected routes(admin only)
	api.Use(middlewares.AdminOnly)
	api.Post("/novels/add-genre",controllers.AddGenre)
	api.Post("/novels/add-tag",controllers.AddTag)
	api.Get("/novels/genre",controllers.ViewGenre)
	api.Get("/novels/tags",controllers.ViewTag)
	api.Post("/novels/add-author",controllers.AddAuthor)
}