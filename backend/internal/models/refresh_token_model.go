package models

import (
	"time"
)

type RefreshToken struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Token      string    `json:"token"`
	UserID     uint      `json:"user_id"`
	ExpiryDate time.Time `json:"expiry_date"`
	User       *User     `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
