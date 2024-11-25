# Chat Application

A real-time chat application built with Go Fiber framework and PostgreSQL.

## Features

- Authentication (JWT)
- Real-time chat with WebSocket
- Chat rooms (Create and join)
- Private messaging
- Online/offline status
- File upload
- Chat history storage
- New message notifications
- Message read status

## Project Structure

chat-app/
├── config/             # Application configuration
├── internal/           # Core source code
│   ├── database/      # DB connection and migration
│   ├── handlers/      # Request handlers
│   ├── middleware/    # Middleware
│   ├── models/        # Data models
│   ├── routes/        # Routing
│   └── websocket/     # WebSocket handling
├── uploads/           # Upload directory
├── main.go            # Entry point
└── README.md

\n\n## Installation




