# Chat Application

Ứng dụng chat realtime được xây dựng với Go Fiber framework và PostgreSQL.

## Tính năng

- Authentication (JWT)
- Real-time chat với WebSocket 
- Phòng chat (Tạo và tham gia)
- Tin nhắn riêng tư
- Trạng thái online/offline
- Upload file
- Lưu trữ lịch sử chat
- Thông báo tin nhắn mới
- Xem trạng thái đã xem tin nhắn

## Cấu trúc Project

chat-app/
├── config/             # Cấu hình ứng dụng
├── internal/           # Mã nguồn chính
│   ├── database/      # Kết nối và migration DB
│   ├── handlers/      # Xử lý request
│   ├── middleware/    # Middleware
│   ├── models/        # Model dữ liệu
│   ├── routes/        # Định tuyến
│   └── websocket/     # Xử lý WebSocket
├── uploads/           # Thư mục lưu file upload
├── main.go           # Entry point
└── README.md

## Cài đặt

1. Clone repository:




