package models

import (
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Users       []User    `json:"users" gorm:"many2many:room_users;"`
	Messages    []Message `json:"messages"`
}

type RoomUser struct {
	RoomID uint `gorm:"primaryKey"`
	UserID uint `gorm:"primaryKey"`
} 