package review

import "github.com/susek555/BD2/car-dealer-api/internal/domains/user"

type Review struct {
	ID          uint       `json:"id" gorm:"primary_key"`
	Description string     `json:"description" gorm:"not null"`
	Rating      int        `json:"rating" gorm:"not null;check:rating BETWEEN 1 AND 5"`
	ReviewerID  uint       `json:"reviewer_id"`
	Reviewer    *user.User `gorm:"foreignKey:ReviewerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReviewedID  uint       `json:"reviewed_id"`
	Reviewed    *user.User `gorm:"foreignKey:ReviewedID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
