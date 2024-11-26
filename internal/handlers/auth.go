package handlers

import (
	"chat-app/internal/constants"
	"chat-app/internal/models"
	"chat-app/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
	//"strings"
)

type AuthHandler struct {
	db     *gorm.DB
	secret string
}

func NewAuthHandler(db *gorm.DB, secret string) *AuthHandler {
	return &AuthHandler{
		db:     db,
		secret: secret,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var registerData struct {
		Username string `json:"username"`
		Password string `json:"password"`
		 Email    string `json:"email"`
	}

	if err := c.BodyParser(&registerData); err != nil {
		log.Printf("Parse error: %v", err)
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"code": constants.StatusBadRequest,
			"message": constants.ErrInvalidRequest,
		})
	}

	// Validate dữ liệu trống
	if registerData.Username == "" || registerData.Password == "" || registerData.Email == "" {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"code": constants.StatusBadRequest,
			"message": "Username, password và email không được để trống",
		})
	}

	// Kiểm tra username đã tồn tại
	var existingUser models.User
	if err := h.db.Where("username = ?", registerData.Username).First(&existingUser).Error; err == nil {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"code": constants.StatusBadRequest,
			"message": "Username đã tồn tại",
		})
	}

	// Kiểm tra email đã tồn tại
	if err := h.db.Where("email = ?", registerData.Email).First(&existingUser).Error; err == nil {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"code": constants.StatusBadRequest,
			"message": "Email đã tồn tại",
		})
	}

	// Debug thông tin
	log.Printf("Register attempt:")
	log.Printf("Username: %s", registerData.Username)
	log.Printf("Password length: %d", len(registerData.Password))
	log.Printf("Email: %s", registerData.Email)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Hash generation failed: %v", err)
		return c.Status(constants.StatusServerError).JSON(fiber.Map{
			"code": constants.StatusServerError,
			"message": constants.ErrHashingPassword,
		})
	}

	// Tạo user mới
	user := models.User{
		Username: registerData.Username,
		Password: string(hashedPassword),
		Email:    registerData.Email,
		LastSeen: time.Now(),
	}

	// Log hash được tạo
	log.Printf("Generated hash: %s", string(hashedPassword))

	// Lưu vào DB
	if err := h.db.Create(&user).Error; err != nil {
		log.Printf("DB save failed: %v", err)
		return c.Status(constants.StatusServerError).JSON(fiber.Map{
			"code": constants.StatusServerError,
			"message": constants.ErrCreatingUser,
		})
	}

	// Verify sau khi lưu
	var savedUser models.User
	h.db.First(&savedUser, user.ID)
	log.Printf("Saved hash: %s", savedUser.Password)

	return c.Status(constants.StatusCreated).JSON(fiber.Map{
		"code": constants.StatusCreated,
		"message": constants.MsgUserCreated,
		"data": fiber.Map{
			"id": user.ID,
			"username": user.Username,
			"email": user.Email,
		},
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(constants.StatusBadRequest).JSON(fiber.Map{
			"code":    constants.StatusBadRequest,
			"message": constants.ErrInvalidRequest,
		})
	}

	var user models.User
	if err := h.db.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
			"code":    constants.StatusUnauthorized,
			"message": constants.ErrInvalidCredentials,
		})
	}

	// Debug thông tin
	log.Printf("Login attempt:")
	log.Printf("Username: %s", loginData.Username)
	log.Printf("Password length: %d", len(loginData.Password))
	log.Printf("Stored hash length: %d", len(user.Password))

	// Kiểm tra xem có ký tự đặc biệt nào không
	log.Printf("Password bytes: %v", []byte(loginData.Password))
	log.Printf("Hash bytes: %v", []byte(user.Password))

	// So sánh password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		log.Printf("Password comparison failed: %v", err)

		// Thử tạo hash mới và so sánh trực tiếp
		newHash, _ := bcrypt.GenerateFromPassword([]byte(loginData.Password), bcrypt.DefaultCost)
		log.Printf("Original hash: %s", user.Password)
		log.Printf("New hash: %s", string(newHash))

		return c.Status(constants.StatusUnauthorized).JSON(fiber.Map{
			"code":    constants.StatusUnauthorized,
			"message": constants.ErrInvalidCredentials,
		})
	}

	token, err := utils.GenerateToken(user.ID, h.secret)
	if err != nil {
		return c.Status(constants.StatusServerError).JSON(fiber.Map{
			"code":    constants.StatusServerError,
			"message": constants.ErrGeneratingToken,
		})
	}

	return c.Status(constants.StatusOK).JSON(fiber.Map{
		"code":    constants.StatusOK,
		"message": constants.MsgLoginSuccess,
		"data": fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		},
	})
}
