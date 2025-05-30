package models

import "time"

type Review struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Description string    `json:"description"`
	ReviewDate  time.Time `json:"date"`
	Rating      uint      `json:"rating"`
	ReviewerID  uint      `json:"reviewer_id"`
	Reviewer    *User     `gorm:"foreignKey:ReviewerID;references:ID"`
	RevieweeId  uint      `json:"reviewee_id"`
	Reviewee    *User     `gorm:"foreignKey:RevieweeId;references:ID"`
}
