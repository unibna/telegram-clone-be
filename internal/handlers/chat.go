package handlers

import (
	"chat-app/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
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

func (h *ChatHandler) SendDirectMessage(c *fiber.Ctx) error {
	senderID := c.Locals("userID").(uint)
	
	var msg struct {
		ReceiverID uint   `json:"receiver_id"`
		Content    string `json:"content"`
	}

	if err := c.BodyParser(&msg); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	// Kiểm tra receiver tồn tại
	var receiver models.User
	if err := h.db.First(&receiver, msg.ReceiverID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	directMsg := models.DirectMessage{
		SenderID:   senderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
	}

	if err := h.db.Create(&directMsg).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error sending message",
		})
	}

	// Load sender info
	h.db.Preload("Sender").First(&directMsg, directMsg.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Message sent successfully",
		"data":    directMsg,
	})
}

func (h *ChatHandler) GetDirectMessages(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	otherUserID := c.QueryInt("user_id")

	var messages []models.DirectMessage
	if err := h.db.Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, otherUserID, otherUserID, userID).
		Order("created_at asc").
		Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching messages",
		})
	}

	// Đánh dấu tin nhắn đã đọc
	h.db.Model(&models.DirectMessage{}).
		Where("receiver_id = ? AND sender_id = ? AND read = ?", userID, otherUserID, false).
		Updates(map[string]interface{}{
			"read":    true,
			"read_at": time.Now(),
		})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Messages fetched successfully",
		"data":    messages,
	})
}