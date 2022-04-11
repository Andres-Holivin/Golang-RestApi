package routes

import (
	"golang-restapi/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	app.Post("/users", controllers.CreateUser)
	app.Get("/user/:userId", controllers.GetAUser)
}
