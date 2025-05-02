package refresh_token

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"time"
)

type RefreshToken struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Token      string    `json:"token" gorm:"unique;not null"`
	UserId     uint      `json:"user_id" gorm:"not null;index"`
	User       user.User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ExpiryDate time.Time `json:"expiry_date" gorm:"column:expiry_date;not null"`
}
