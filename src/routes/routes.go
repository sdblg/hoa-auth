package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sdblg/hoa-auth/src/controllers"
)

func Setup(app *fiber.App) {
	app.Get("/", controllers.HealthCheck)

	app.Post("/api/v1/user/register", controllers.UserRegister)
	app.Post("/api/v1/user/login", controllers.UserLogin)
	app.Post("/api/v1/user/logout", controllers.UserLogout)
	app.Get("/api/v1/user", controllers.User)
}
