package handlers

import (
	"chat-app/internal/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

type ContactResponse struct {
	ID       uint   `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	IsOnline bool   `json:"is_online"`
	LastSeen string `json:"last_seen"`
}

func (h *UserHandler) ListUser(c *fiber.Ctx) error {
	var users []models.User
	id := c.Query("id")
	username := c.Query("username")
	email := c.Query("email")

	query := h.db

	if id != "" {
		query = query.Where("id = ?", id)
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	if err := query.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching users",
		})
	}

	if len(users) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No users found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data":   users,
	})
}

func (h *UserHandler) AddContactUser(c *fiber.Ctx) error {
	var contact models.Contact
	if err := c.BodyParser(&contact); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("Error parsing request: %v", err),
		})
	}

	// Check if the user with UserID exists
	var user models.User
	if err := h.db.First(&user, contact.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": fmt.Sprintf("User not found: %v", contact.UserID),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error checking user existence",
		})
	}
	var userContact models.User
	// Check if the user with UserContactID exists
	if err := h.db.First(&userContact, contact.UserContactID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": fmt.Sprintf("User not found: %v", contact.UserContactID),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error checking contact user existence",
		})
	}

	// Create the contact
	if err := h.db.Create(&contact).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Error creating contact: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data":   contact,
	})
}

func (h *UserHandler) ListContactUser(c *fiber.Ctx) error {
	var contacts []models.Contact
	userID := c.Query("user_id")

	query := h.db

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Preload the ContactUser field
	if err := query.Preload("ContactUser").Find(&contacts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Error fetching contacts: %v", err),
		})
	}

	if len(contacts) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("No contacts found"),
		})
	}

	var contactResponses []ContactResponse
	for _, contact := range contacts {
		contactResponses = append(contactResponses, ContactResponse{
			ID:       contact.UserContactID,
			UserName: contact.ContactUser.Username,
			Email:    contact.ContactUser.Email,
			IsOnline: contact.ContactUser.IsOnline,
			LastSeen: contact.ContactUser.LastSeen.String(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data":   contactResponses,
	})
}

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	var user models.User
	userID := c.Locals("userID").(uint)

	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"data":    user,
	})
}
