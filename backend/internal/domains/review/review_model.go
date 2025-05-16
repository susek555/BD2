package review

import "github.com/susek555/BD2/car-dealer-api/internal/domains/user"

type Review struct {
	ID          uint       `json:"id" gorm:"primary_key;autoIncrement"`
	Description string     `json:"description" gorm:"not null"`
	ReviewDate  string     `json:"date" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Rating      uint       `json:"rating" gorm:"not null;check:rating BETWEEN 1 AND 5"`
	ReviewerID  uint       `json:"reviewer_id" gorm:"not null;check:reviewer_id <> reviewee_id;uniqueIndex:ux_review_pair"`
	Reviewer    *user.User `gorm:"foreignKey:ReviewerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RevieweeId  uint       `json:"reviewee_id" gorm:"not null;uniqueIndex:ux_review_pair"`
	Reviewee    *user.User `gorm:"foreignKey:RevieweeId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
