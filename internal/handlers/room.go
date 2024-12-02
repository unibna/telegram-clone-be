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
		"status":  fiber.StatusCreated,
		"message": "Room created successfully",
		"data":    room,
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
		"status":  fiber.StatusOK,
		"message": "Successfully joined room",
	})
}

// Get all of the rooms of user
func (h *RoomHandler) GetMyChatRooms(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	roomUsers := []models.RoomUser{}
	if err := h.db.Model(&models.RoomUser{}).
		Where("user_id = ?", userID).
		Find(&roomUsers).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	roomIDs := []uint{}

	for _, roomUser := range roomUsers {
		roomIDs = append(roomIDs, roomUser.RoomID)
	}

	rooms := []models.Room{}

	if err := h.db.Model(&models.Room{}).
		Where("id IN ?", roomIDs).
		Find(&rooms).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Get user rooms successfully",
		"data":    rooms,
	})
}

// Get all messages inside a chat room.
func (h *RoomHandler) GetRoomMessages(c *fiber.Ctx) error {
	roomID := c.Params("roomID")

	messages := []models.Message{}

	if err := h.db.Model(&models.Message{}).
		Where("room_id = ?", roomID).
		Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    messages,
		"message": "Get room message successfully",
	})
}
