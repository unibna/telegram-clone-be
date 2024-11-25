package handlers

import (
	"chat-app/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ChatHandler struct {
	db *gorm.DB
}

func NewChatHandler(db *gorm.DB) *ChatHandler {
	return &ChatHandler{
		db: db,
	}
}

func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	
	var message models.Message
	if err := c.BodyParser(&message); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	message.UserID = userID

	if err := h.db.Create(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error sending message",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": fiber.StatusCreated,
		"message": "Message sent successfully",
		"data": message,
	})
}

func (h *ChatHandler) GetMessages(c *fiber.Ctx) error {
	var messages []models.Message
	if err := h.db.Preload("User").Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching messages",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data": messages,
	})
} 