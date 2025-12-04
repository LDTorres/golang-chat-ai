package v1

import (
	"github.com/LDTorres/golang-chat-ai/internal/database"
	"github.com/LDTorres/golang-chat-ai/internal/integrations/llm"
	"github.com/LDTorres/golang-chat-ai/internal/models"
	"github.com/gofiber/fiber/v2"
)

// Initialize LLM provider
var llmProvider llm.LLMProvider

func InitLLM() {
	var err error
	llmProvider, err = llm.NewLLMProvider()
	if err != nil {
		// Fallback to mock if config fails or not set, or handle error
		// For now, let's just log and use mock if it fails, or maybe panic?
		// Given the requirements, let's try to be robust.
		llmProvider = &llm.MockLLM{}
	}
}

func incrementMessageCount(userID uint) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err == nil {
		user.MessageCount++
		database.DB.Save(&user)
	}
}

func CreateChat(c *fiber.Ctx) error {
	type Request struct {
		UserID  uint   `json:"user_id"`
		Message string `json:"message"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if len(req.Message) > 300 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Message exceeds 300 characters"})
	}

	// Create Chat
	chat := models.Chat{
		UserID: req.UserID,
		Title:  req.Message, // Use first message as title for now
	}
	if err := database.DB.Create(&chat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create chat"})
	}

	// Save User Message
	userMsg := models.Message{
		ChatID:  chat.ID,
		Role:    "user",
		Content: req.Message,
	}
	database.DB.Create(&userMsg)
	incrementMessageCount(req.UserID)

	prompt := req.Message

	// Get LLM Response
	response, _, err := llmProvider.GenerateResponse(prompt, "")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate response"})
	}

	// Save Assistant Message
	assistantMsg := models.Message{
		ChatID:  chat.ID,
		Role:    "assistant",
		Content: response,
	}
	database.DB.Create(&assistantMsg)

	return c.JSON(fiber.Map{
		"chat":     chat,
		"response": assistantMsg,
	})
}

func DeleteChat(c *fiber.Ctx) error {
	id := c.Params("id")
	// Soft delete
	if err := database.DB.Delete(&models.Chat{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete chat"})
	}
	return c.SendStatus(fiber.StatusOK)
}

func GetChats(c *fiber.Ctx) error {
	userID := c.Params("id")
	var chats []models.Chat
	if err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&chats).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch chats"})
	}
	return c.JSON(chats)
}

func GetMessages(c *fiber.Ctx) error {
	chatID := c.Params("id")
	var messages []models.Message
	if err := database.DB.Where("chat_id = ?", chatID).Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch messages"})
	}
	return c.JSON(messages)
}

func SendMessage(c *fiber.Ctx) error {
	chatID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid chat ID"})
	}

	type Request struct {
		Message string `json:"message"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if len(req.Message) > 300 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Message exceeds 300 characters"})
	}

	// We need UserID to check limits. Fetch chat first.
	var chat models.Chat
	if err := database.DB.First(&chat, chatID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Chat not found"})
	}

	var previousAssistantMessage models.Message
	database.DB.Last(&previousAssistantMessage, map[string]interface{}{
		"chat_id": chatID,
		"role":    "assistant",
	})

	// Save User Message
	userMsg := models.Message{
		ChatID:  uint(chatID),
		Role:    "user",
		Content: req.Message,
	}
	database.DB.Create(&userMsg)
	incrementMessageCount(chat.UserID)

	prompt := req.Message

	var previousId string = previousAssistantMessage.ModelMessageId

	// Get LLM Response
	response, id, err := llmProvider.GenerateResponse(prompt, previousId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate response"})
	}

	// Save Assistant Message
	assistantMsg := models.Message{
		ChatID:         uint(chatID),
		Role:           "assistant",
		Content:        response,
		ModelMessageId: id,
	}
	database.DB.Create(&assistantMsg)

	return c.JSON(assistantMsg)
}

func Chats(app fiber.Router) {
	api := app.Group("/chats")
	api.Post("/", CreateChat)
	api.Get("/:id/messages", GetMessages)
	api.Post("/:id/messages", SendMessage)
}
