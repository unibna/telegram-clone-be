package handlers

import (
	"chat-app/internal/models"
	"chat-app/pkg/utils"
	ws "chat-app/internal/websocket"
	"encoding/json"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type WebSocketHandler struct {
	db     *gorm.DB
	Hub    *ws.Hub
	secret string
}

func NewWebSocketHandler(db *gorm.DB, secret string) *WebSocketHandler {
	return &WebSocketHandler{
		db:     db,
		Hub:    ws.NewHub(),
		secret: secret,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *websocket.Conn) {
	log.Println("New WebSocket connection attempt...")
	
	// Extract token from Authorization header
	authHeader := c.Headers("Authorization")
	if authHeader == "" {
		log.Println("Missing Authorization header")
		return
	}

	// Remove Bearer prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	
	// Parse and validate token
	userID, err := utils.ValidateToken(tokenString, h.secret)
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return
	}

	log.Printf("User %d authenticated successfully", userID)

	// Create done channel to coordinate shutdown
	done := make(chan struct{})

	client := &ws.Client{
		Hub:     h.Hub,
		Conn:    c,
		Send:    make(chan []byte, 256),
		UserID:  userID,
		IsAlive: true,
	}

	// Start handlers before registering with hub
	go func() {
		h.writePump(client)
		close(done)
	}()
	
	go h.readPump(client, done)

	log.Printf("WebSocket handlers started for user %d", userID)

	// Update user status
	if err := h.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"is_online":  true,
		"last_seen": time.Now(),
	}).Error; err != nil {
		log.Printf("Error updating user status: %v", err)
	}

	// Register with hub after handlers are running
	h.Hub.Register <- client

	// Keep the connection alive until done channel is closed
	<-done
	log.Printf("WebSocket connection closed for user %d", userID)
}

func (h *WebSocketHandler) readPump(client *ws.Client, done chan struct{}) {
	defer func() {
		log.Printf("ReadPump closing for user %d", client.UserID)
		h.Hub.Unregister <- client
		client.Conn.Close()
		// Update user status when disconnected
		if err := h.db.Model(&models.User{}).Where("id = ?", client.UserID).Updates(map[string]interface{}{
			"is_online":  false,
			"last_seen": time.Now(),
		}).Error; err != nil {
			log.Printf("Error updating user offline status: %v", err)
		}
	}()

	client.Conn.SetReadLimit(ws.MaxMessageSize)
	if err := client.Conn.SetReadDeadline(time.Now().Add(ws.PongWait)); err != nil {
		log.Printf("Error setting read deadline for user %d: %v", client.UserID, err)
		return
	}
	
	client.Conn.SetPongHandler(func(string) error {
		err := client.Conn.SetReadDeadline(time.Now().Add(ws.PongWait))
		if err != nil {
			log.Printf("Error updating read deadline in pong handler for user %d: %v", client.UserID, err)
			return err
		}
		return nil
	})

	for {
		select {
		case <-done:
			return
		default:
			messageType, message, err := client.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("ReadPump error for user %d: %v", client.UserID, err)
				} else {
					log.Printf("Normal close for user %d: %v", client.UserID, err)
				}
				return
			}
			
			log.Printf("Received message type %d from user %d", messageType, client.UserID)

			// Handle ping message
			if messageType == websocket.PingMessage {
				if err := client.Conn.WriteMessage(websocket.PongMessage, nil); err != nil {
					log.Printf("Error sending pong to user %d: %v", client.UserID, err)
					return
				}
				continue
			}

			// Only process text messages
			if messageType != websocket.TextMessage {
				log.Printf("Ignoring non-text message type %d from user %d", messageType, client.UserID)
				continue
			}

			// Parse message
			var msg ws.Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Printf("Error parsing message from user %d: %v", client.UserID, err)
				continue
			}

			// Handle direct messages
			if msg.Type == "direct" {
				// Validate receiver exists
				var receiver models.User
				if err := h.db.First(&receiver, msg.ReceiverID).Error; err != nil {
					log.Printf("Receiver %d not found for message from user %d", msg.ReceiverID, client.UserID)
					notification := ws.Message{
						Type:    "error",
						Content: "Receiver not found",
					}
					notificationJSON, _ := json.Marshal(notification)
					client.Send <- notificationJSON
					continue
				}

				// Save message to database
				directMsg := models.DirectMessage{
					SenderID:   client.UserID,
					ReceiverID: msg.ReceiverID,
					Content:    msg.Content,
				}

				if err := h.db.Create(&directMsg).Error; err != nil {
					log.Printf("Error saving message from user %d: %v", client.UserID, err)
					notification := ws.Message{
						Type:    "error",
						Content: "Failed to save message",
					}
					notificationJSON, _ := json.Marshal(notification)
					client.Send <- notificationJSON
					continue
				}

				log.Printf("Message saved to database: ID=%d, From=%d, To=%d", directMsg.ID, client.UserID, msg.ReceiverID)

				// Prepare message for sending
				msg.MessageID = directMsg.ID
				msg.SenderID = client.UserID
				msg.CreatedAt = directMsg.CreatedAt

				// Send to receiver if online
				messageJSON, _ := json.Marshal(msg)
				delivered := h.Hub.SendToUser(msg.ReceiverID, messageJSON)

				// Update delivered status
				if delivered {
					if err := h.db.Model(&directMsg).Update("delivered", true).Error; err != nil {
						log.Printf("Error updating message delivery status: %v", err)
					}
					msg.Delivered = true
					log.Printf("Message delivered to user %d", msg.ReceiverID)
				}

				// Send confirmation to sender
				client.Send <- messageJSON

				// If receiver is offline, send notification to sender
				if !delivered {
					log.Printf("Receiver %d is offline, sending notification to sender %d", msg.ReceiverID, client.UserID)
					notification := ws.Message{
						Type:    "notification",
						Content: "User is offline. Message will be delivered when they come online.",
					}
					notificationJSON, _ := json.Marshal(notification)
					client.Send <- notificationJSON
				}
			}
		}
	}
}

func (h *WebSocketHandler) writePump(client *ws.Client) {
	ticker := time.NewTicker(ws.PingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
		log.Printf("WritePump closing for user %d", client.UserID)
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(ws.WriteWait))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Error getting next writer for user %d: %v", client.UserID, err)
				return
			}
			w.Write(message)

			// Add queued messages
			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				log.Printf("Error closing writer for user %d: %v", client.UserID, err)
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(ws.WriteWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending ping to user %d: %v", client.UserID, err)
				return
			}
		}
	}
}