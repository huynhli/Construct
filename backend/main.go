package main

import (
	"backend/config"
	"backend/database"
	"backend/routes"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func init() {
	config.LoadConfig()
}

func main() {
	fmt.Println("Hello, Go project!")

	database.ConnectPostgres()
	defer database.DisconnectPostgres()

	if config.Port == "" {
		config.Port = "8080"
	}

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		return c.Next()
	})

	config.SetupCors(app)
	routes.SetupRoutes(app)

	for _, route := range app.GetRoutes() {
		fmt.Printf("Method: %s, Path: %s\n", route.Method, route.Path)
	}

	err := app.Listen(":" + config.Port)
	if err != nil {
		log.Fatalf("Error listening to app on port")
	}
}
