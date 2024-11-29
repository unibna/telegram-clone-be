package main

import (
	"chat-app/config"
	"chat-app/internal/database"
	"chat-app/internal/handlers"
	"chat-app/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
	h := &routes.Handlers{
		AuthHandler: handlers.NewAuthHandler(db, cfg.JWTSecret),
		ChatHandler: handlers.NewChatHandler(db),
		RoomHandler: handlers.NewRoomHandler(db),
		UserHandler: handlers.NewUserHandler(db),
		WSHandler:   handlers.NewWebSocketHandler(db, cfg.JWTSecret),
	}

	// Khởi chạy WebSocket hub
	go h.WSHandler.Hub.Run()

	// Setup routes
	routes.SetupRoutes(app, h, cfg.JWTSecret)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
