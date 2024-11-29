package routes

import (
	"chat-app/internal/handlers"
	"chat-app/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Handlers struct {
	AuthHandler *handlers.AuthHandler
	ChatHandler *handlers.ChatHandler
	RoomHandler *handlers.RoomHandler
	UserHandler *handlers.UserHandler
	WSHandler   *handlers.WebSocketHandler
}

func SetupRoutes(app *fiber.App, h *Handlers, jwtSecret string) {
	// WebSocket route
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Use("/ws", middleware.AuthMiddleware(jwtSecret))
	app.Get("/ws", websocket.New(h.WSHandler.HandleWebSocket))

	// API routes
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", h.AuthHandler.Register)
	auth.Post("/login", h.AuthHandler.Login)

	// Protected routes
	protected := api.Group("/", middleware.AuthMiddleware(jwtSecret))

	// Chat routes
	chat := protected.Group("/chat")
	chat.Post("/messages", h.ChatHandler.SendMessage)
	chat.Get("/messages", h.ChatHandler.GetMessages)
	chat.Post("/direct", h.ChatHandler.SendDirectMessage)
	chat.Get("/direct", h.ChatHandler.GetDirectMessages)

	// Room routes
	rooms := protected.Group("/rooms")
	rooms.Post("/", h.RoomHandler.CreateRoom)
	rooms.Post("/:roomID/join", h.RoomHandler.JoinRoom)
	rooms.Get("/me", h.RoomHandler.GetMyChatRooms)
	rooms.Get("/:roomID/messages", h.RoomHandler.GetRoomMessages)

	// User routes
	users := protected.Group("/users")
	users.Get("/", h.UserHandler.ListUser)
	users.Post("/add_contact", h.UserHandler.AddContactUser)
	users.Get("/list_contact", h.UserHandler.ListContactUser)
}
