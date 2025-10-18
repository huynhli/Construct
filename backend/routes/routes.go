package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	app.Get("/", homePage)
}

func homePage(c *fiber.Ctx) error {
	ret := "hey this works"
	return c.SendString(ret)
}
