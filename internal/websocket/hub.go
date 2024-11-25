package websocket

import (
	"encoding/json"
	"github.com/gofiber/websocket/v2"
	"log"
	"sync"
	"time"
)

type Message struct {
	Type      string  `json:"type"`
	Content   string  `json:"content"`
	UserID    uint    `json:"user_id"`
	RoomID    string  `json:"room_id"`
	ToUserID  *uint   `json:"to_user_id,omitempty"`
	FileURL   string  `json:"file_url,omitempty"`
	Timestamp string  `json:"timestamp"`
}

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID uint
	RoomID string
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	rooms      map[string]map[*Client]bool
	userStatus map[uint]bool
	mu         sync.Mutex
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.Hub.broadcast <- message
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
		userStatus: make(map[uint]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
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

		case client := <-h.Unregister:
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