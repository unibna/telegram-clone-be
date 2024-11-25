package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username    string    `json:"username" gorm:"unique"`
	Password    string    `json:"-"`
	Email       string    `json:"email" gorm:"unique"`
	IsOnline    bool      `json:"is_online" gorm:"default:false"`
	LastSeen    time.Time `json:"last_seen"`
	Rooms       []Room    `json:"rooms" gorm:"many2many:room_users;"`
} 