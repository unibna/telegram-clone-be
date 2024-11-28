package websocket

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

const (
	// Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * 9) / 10

	// Maximum message size allowed from peer.
	MaxMessageSize = 512 * 1024 // 512KB
)

type Message struct {
	Type       string    `json:"type"` // "direct", "notification"
	MessageID  uint      `json:"message_id,omitempty"`
	SenderID   uint      `json:"sender_id,omitempty"`
	ReceiverID uint      `json:"receiver_id,omitempty"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	Delivered  bool      `json:"delivered,omitempty"`
	Read       bool      `json:"read,omitempty"`
}

type Client struct {
	Hub     *Hub
	Conn    *websocket.Conn
	Send    chan []byte
	UserID  uint
	RoomID  string
	mu      sync.Mutex
	IsAlive bool
}

type Hub struct {
	mu         sync.RWMutex
	clients    map[*Client]bool
	broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	UserConns  map[uint]*Client // Map userID to their connection
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		UserConns:  make(map[uint]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			// Unregister existing connection if any
			if oldClient, exists := h.UserConns[client.UserID]; exists {
				delete(h.clients, oldClient)
				close(oldClient.Send)
			}
			h.clients[client] = true
			h.UserConns[client.UserID] = client
			client.IsAlive = true
			h.mu.Unlock()
			log.Printf("Client connected. UserID: %d", client.UserID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.UserConns, client.UserID)
				close(client.Send)
				client.IsAlive = false
			}
			h.mu.Unlock()
			log.Printf("Client disconnected. UserID: %d", client.UserID)
		}
	}
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(targetUserID uint, message []byte) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.UserConns[targetUserID]; ok && client.IsAlive {
		select {
		case client.Send <- message:
			return true
		default:
			return false
		}
	}
	return false
}

// IsUserOnline checks if a user is currently connected
func (h *Hub) IsUserOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.UserConns[userID]; ok {
		return client.IsAlive
	}
	return false
}
