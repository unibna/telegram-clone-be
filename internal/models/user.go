package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string    `json:"username" gorm:"unique;not null"`
	Password string    `json:"-" gorm:"not null"`
	Email    string    `json:"email" gorm:"unique;not null"`
	IsOnline bool      `json:"is_online" gorm:"default:false"`
	LastSeen time.Time `json:"last_seen"`
	Rooms    []Room    `json:"rooms" gorm:"many2many:room_users;"`
}
