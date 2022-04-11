package main

import (
	"golang-restapi/configs"
	"golang-restapi/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	configs.ConnectDB()

	routes.UserRoute(app)

	app.Listen(":6000")
}
