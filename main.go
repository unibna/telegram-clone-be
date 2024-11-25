package main

import (
	"chat-app/config"
	"chat-app/internal/database"
	"chat-app/internal/handlers"
	"chat-app/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"log"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
	chatHandler := handlers.NewChatHandler(db)
	roomHandler := handlers.NewRoomHandler(db)
	wsHandler := handlers.NewWebSocketHandler(db)

	// Khởi chạy WebSocket hub
	go wsHandler.hub.Run()

	// WebSocket route
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(wsHandler.HandleWebSocket))

	// Routes
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Chat routes (protected)
	chat := api.Group("/chat", middleware.AuthMiddleware(cfg.JWTSecret))
	chat.Post("/messages", chatHandler.SendMessage)
	chat.Get("/messages", chatHandler.GetMessages)

	// Room routes
	rooms := api.Group("/rooms", middleware.AuthMiddleware(cfg.JWTSecret))
	rooms.Post("/", roomHandler.CreateRoom)
	rooms.Post("/:roomID/join", roomHandler.JoinRoom)

	// Start server
	log.Fatal(app.Listen(":3000"))
} 