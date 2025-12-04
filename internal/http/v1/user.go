package v1

import (
	"github.com/LDTorres/golang-chat-ai/internal/database"
	"github.com/LDTorres/golang-chat-ai/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateUser(c *fiber.Ctx) error {
	type Request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	// Check if user exists
	var user models.User
	result := database.DB.Where("email = ?", req.Email).First(&user)
	if result.Error == nil {
		// User found, return it
		return c.Status(fiber.StatusOK).JSON(user)
	}

	// User not found, create new
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name is required for new users",
		})
	}

	user = models.User{
		Name:     req.Name,
		Email:    req.Email,
		PublicID: uuid.New().String(),
	}

	if result := database.DB.Create(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}
