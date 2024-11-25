package handlers

import (
	"chat-app/internal/models"
	ws "chat-app/internal/websocket"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
	"time"
)

type WebSocketHandler struct {
	db  *gorm.DB
	Hub *ws.Hub
}

func NewWebSocketHandler(db *gorm.DB) *WebSocketHandler {
	return &WebSocketHandler{
		db:  db,
		Hub: ws.NewHub(),
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *websocket.Conn) {
	userID := c.Locals("userID").(uint)
	roomID := c.Query("room_id")

	client := &ws.Client{
		Hub:    h.Hub,
		Conn:   c,
		Send:   make(chan []byte, 256),
		UserID: userID,
		RoomID: roomID,
	}

	h.Hub.Register <- client

	// Update user status
	h.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"is_online":  true,
		"last_seen": time.Now(),
	})

	go client.WritePump()
	go client.ReadPump()
} 