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
	DB_HOST     string
	DB_NAME     string
	DB_PORT     string
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
	DB_HOST = os.Getenv("DB_HOST")
	DB_NAME = os.Getenv("DB_NAME")
	DB_PORT = os.Getenv("DB_PORT")

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
