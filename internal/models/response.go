package models

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func NewResponse(code int, message string, data interface{}) *Response {
    return &Response{
        Code:    code,
        Message: message,
        Data:    data,
    }
}