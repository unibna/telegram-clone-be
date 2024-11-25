package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Content    string `json:"content"`
	UserID     uint   `json:"user_id"`
	User       User   `json:"user" gorm:"foreignKey:UserID"`
	RoomID     uint   `json:"room_id"`
	Room       Room   `json:"room" gorm:"foreignKey:RoomID"`
	IsPrivate  bool   `json:"is_private"`
	ToUserID   *uint  `json:"to_user_id,omitempty"`
	ToUser     *User  `json:"to_user,omitempty" gorm:"foreignKey:ToUserID"`
	FileURL    string `json:"file_url,omitempty"`
} 