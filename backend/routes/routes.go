package routes

import (
	"backend/handlers"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	ver := api.Group("/v1")

	ver.Get("/", homePage)
	ver.Post("/signup", handlers.Signup)
	ver.Post("/login", handlers.Login)
	ver.Get("/tasks", JWTMiddleware(), handlers.UserGetTasks)
	ver.Get("/projects", JWTMiddleware(), handlers.UserGetProjects)
	ver.Delete("/projects", JWTMiddleware(), handlers.AdminDeleteProject)
	ver.Post("/projects", JWTMiddleware(), handlers.AdminAddOrEditProject)
	ver.Delete("/tasks", JWTMiddleware(), handlers.DeleteTask)
	ver.Post("/tasks", JWTMiddleware(), handlers.AdminAddOrEditTask)
}

func homePage(c *fiber.Ctx) error {
	ret := "hey this works"
	return c.SendString(ret)
}

// auth middleware to ensure authorized user
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer"))
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
		}

		username, ok := claims["username"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
		}

		isAdmin, ok := claims["isAdmin"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
		}

		userID := int(claims["id"].(float64)) // note JSON numbers are float64

		c.Locals("username", username)
		c.Locals("userID", userID)
		c.Locals("isAdmin", isAdmin)

		return c.Next()
	}
}
