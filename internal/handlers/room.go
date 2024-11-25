package handlers

import (
	"chat-app/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RoomHandler struct {
	db *gorm.DB
}

func NewRoomHandler(db *gorm.DB) *RoomHandler {
	return &RoomHandler{
		db: db,
	}
}

func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	var room models.Room
	if err := c.BodyParser(&room); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	if err := h.db.Create(&room).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating room",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": fiber.StatusCreated,
		"message": "Room created successfully",
		"data": room,
	})
}

func (h *RoomHandler) JoinRoom(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	roomID := c.Params("roomID")

	var user models.User
	var room models.Room

	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := h.db.First(&room, roomID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Room not found",
		})
	}

	if err := h.db.Model(&room).Association("Users").Append(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error joining room",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"message": "Successfully joined room",
	})
} 