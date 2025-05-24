package models

import (
	"time"
)

type RefreshToken struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Token      string    `json:"token" gorm:"unique;not null"`
	UserId     uint      `json:"user_id" gorm:"not null;index"`
	User       User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ExpiryDate time.Time `json:"expiry_date" gorm:"type:timestamptz;column:expiry_date;not null"`
}
