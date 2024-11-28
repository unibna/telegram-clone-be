package models

import (
	"gorm.io/gorm"
	"time"
)

type DirectMessage struct {
	gorm.Model
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	Content    string    `json:"content"`
	Read       bool      `json:"read" gorm:"default:false"`
	ReadAt     time.Time `json:"read_at,omitempty"`
	Delivered  bool      `json:"delivered" gorm:"default:false"`
	Sender     User      `json:"sender" gorm:"foreignKey:SenderID"`
	Receiver   User      `json:"receiver" gorm:"foreignKey:ReceiverID"`
}
