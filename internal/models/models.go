package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"uniqueIndex"`
	PublicID     string `json:"public_id" gorm:"uniqueIndex"`
	MessageCount int    `json:"message_count" gorm:"default:0"`
	Chats        []Chat `json:"chats"`
}

type Chat struct {
	gorm.Model
	UserID   uint      `json:"user_id"`
	Title    string    `json:"title"` // Optional: First message or summary
	Messages []Message `json:"messages"`
}

type Message struct {
	gorm.Model
	ChatID         uint   `json:"chat_id"`
	Role           string `json:"role"` // "user" or "assistant"
	Content        string `json:"content"`
	ModelMessageId string `json:"model_message_id" gorm:"default:null"`
}
