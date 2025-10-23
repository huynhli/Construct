package routes

import (
	"backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	ver := api.Group("/v1")

	ver.Get("/", homePage)

	ver.Get("/tasks", handlers.UserGetTasks)
	ver.Get("/projects", handlers.UserGetProjects)
}

func homePage(c *fiber.Ctx) error {
	ret := "hey this works"
	return c.SendString(ret)
}
