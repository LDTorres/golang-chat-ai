package main

import (
	"log"
	"os"

	"github.com/LDTorres/golang-chat-ai/internal/database"
	v1 "github.com/LDTorres/golang-chat-ai/internal/http/v1"
	"github.com/LDTorres/golang-chat-ai/internal/models"
	"github.com/LDTorres/golang-chat-ai/internal/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Config .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database
	database.Connect()
	database.DB.AutoMigrate(&models.User{}, &models.Chat{}, &models.Message{})

	// Create a new engine
	engine := mustache.New("./views", ".mustache")

	// Create app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Logger
	shared.Logger(app)

	// Rate Limit
	shared.RateLimit(app)

	// Health Check
	shared.HealthCheck(app)

	// Statics
	app.Static("/static", "./public")

	// View engine routes
	app.Get("/", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("index", fiber.Map{
			"Title": "Login",
		}, "layouts/main")
	})

	app.Get("/chat", func(c *fiber.Ctx) error {
		return c.Render("chat", fiber.Map{
			"Title": "Chat",
		}, "layouts/main")
	})

	// API routes
	v1.ApiV1(app)

	HOST := os.Getenv("HOST")
	PORT := os.Getenv("PORT")
	log.Fatal(app.Listen(HOST + ":" + PORT))
}
