package models

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	UserID        uint      `json:"user_id"`
	UserContactID uint      `json:"user_contact_id"`
	User          *User     `gorm:"foreignKey:UserID;references:ID"`
	ContactUser   *User     `gorm:"foreignKey:UserContactID;references:ID"`
	LastSeen      time.Time `json:"last_seen"`
}
