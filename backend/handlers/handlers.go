package handlers

import "github.com/gofiber/fiber/v2"

func GetTasks(c *fiber.Ctx) error {
	return c.SendString("T")
}

func GetProjects(c *fiber.Ctx) error {
	return c.SendString("P")
}
