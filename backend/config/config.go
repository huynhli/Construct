package config

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var (
	DB_USERNAME string
	DB_PASSWORD string
	Port        string
)

func LoadConfig() {
	if os.Getenv("DOCKER") != "true" {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	// Set vars using os.Getenv("KEY")
	DB_USERNAME = os.Getenv("DB_USERNAME")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	Port = os.Getenv("PORT")
}

func SetupCors(app *fiber.App) {
	corsSetup := cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	})

	app.Use(corsSetup)
}
