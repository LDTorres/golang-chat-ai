package v1

import "github.com/gofiber/fiber/v2"

func ApiV1(app *fiber.App) {
	// Init LLM
	InitLLM()

	v1 := app.Group("/api/v1")

	// Users
	v1.Post("/users", CreateUser)
	v1.Get("/users/:id/chats", GetChats)

	// Chats
	v1.Delete("/chats/:id", DeleteChat) // Register DeleteChat route
	Chats(v1)
}
