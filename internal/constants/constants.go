package constants

import "github.com/gofiber/fiber/v2"

const (
    // HTTP Status codes
    StatusOK           = fiber.StatusOK
    StatusCreated      = fiber.StatusCreated 
    StatusBadRequest   = fiber.StatusBadRequest
    StatusUnauthorized = fiber.StatusUnauthorized
    StatusForbidden    = fiber.StatusForbidden
    StatusNotFound     = fiber.StatusNotFound
    StatusServerError  = fiber.StatusInternalServerError

    // Error messages
    ErrInvalidRequest      = "Yêu cầu không hợp lệ"
    ErrHashingPassword     = "Lỗi mã hóa mật khẩu"
    ErrCreatingUser        = "Lỗi tạo người dùng"
    ErrInvalidCredentials  = "Thông tin đăng nhập không chính xác"
    ErrGeneratingToken     = "Lỗi tạo token"
    ErrUnauthorized        = "Chưa xác thực"
    ErrTokenInvalid        = "Token không hợp lệ"
    ErrUserNotFound        = "Không tìm thấy người dùng"
    ErrRoomNotFound        = "Không tìm thấy phòng chat"
    ErrServerError         = "Lỗi hệ thống"
    ErrDatabaseError       = "Lỗi cơ sở dữ liệu"
    ErrWebSocketError      = "Lỗi kết nối WebSocket"
    
    // Chat errors
    ErrSendingMessage      = "Lỗi gửi tin nhắn"
    ErrFetchingMessages    = "Lỗi lấy tin nhắn"

    // Room errors  
    ErrCreatingRoom        = "Lỗi tạo phòng"
    ErrJoiningRoom         = "Lỗi tham gia phòng"

    // Upload errors
    ErrFileUploadError     = "Lỗi tải file"
    ErrGettingFile         = "Lỗi lấy file"
    ErrSavingFile          = "Lỗi lưu file"

    // Success messages
    MsgUserCreated         = "Tạo tài khoản thành công"
    MsgLoginSuccess        = "Đăng nhập thành công"
    MsgMessageSent         = "Gửi tin nhắn thành công"
    MsgMessagesFetched     = "Lấy tin nhắn thành công"
    MsgRoomCreated         = "Tạo phòng thành công"
    MsgRoomJoined          = "Tham gia phòng thành công"
    MsgFileUploaded        = "Tải file thành công"
)