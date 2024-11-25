package websocket

import (
	"chat-app/internal/models"
	"encoding/json"
	"sync"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   uint
	RoomID   uint
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	rooms      map[uint]map[*Client]bool
	userStatus map[uint]bool
	mu         sync.Mutex
}

type Message struct {
	Type      string          `json:"type"`
	Content   string          `json:"content"`
	UserID    uint           `json:"user_id"`
	RoomID    uint           `json:"room_id"`
	ToUserID  *uint          `json:"to_user_id,omitempty"`
	FileURL   string         `json:"file_url,omitempty"`
	Timestamp string         `json:"timestamp"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[uint]map[*Client]bool),
		userStatus: make(map[uint]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if _, ok := h.rooms[client.RoomID]; !ok {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			h.userStatus[client.UserID] = true
			h.mu.Unlock()

			// Broadcast online status
			status := Message{
				Type:    "status",
				UserID:  client.UserID,
				Content: "online",
			}
			statusJSON, _ := json.Marshal(status)
			h.broadcast <- statusJSON

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.rooms[client.RoomID], client)
				h.userStatus[client.UserID] = false
				close(client.Send)
			}
			h.mu.Unlock()

			// Broadcast offline status
			status := Message{
				Type:    "status",
				UserID:  client.UserID,
				Content: "offline",
			}
			statusJSON, _ := json.Marshal(status)
			h.broadcast <- statusJSON

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					h.mu.Lock()
					delete(h.clients, client)
					delete(h.rooms[client.RoomID], client)
					close(client.Send)
					h.mu.Unlock()
				}
			}
		}
	}
} 