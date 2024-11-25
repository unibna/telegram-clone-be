package handlers

import (
	"chat-app/internal/models"
	"chat-app/internal/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
	"time"
)

type WebSocketHandler struct {
	db  *gorm.DB
	hub *websocket.Hub
}

func NewWebSocketHandler(db *gorm.DB) *WebSocketHandler {
	return &WebSocketHandler{
		db:  db,
		hub: websocket.NewHub(),
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *websocket.Conn) {
	userID := c.Locals("userID").(uint)
	roomID := c.Query("room_id")

	client := &websocket.Client{
		Hub:    h.hub,
		Conn:   c,
		Send:   make(chan []byte, 256),
		UserID: userID,
		RoomID: roomID,
	}

	client.Hub.Register <- client

	// Update user status
	h.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"is_online":  true,
		"last_seen": time.Now(),
	})

	go client.WritePump()
	go client.ReadPump()
} 