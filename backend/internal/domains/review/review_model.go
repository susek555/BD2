package review

import "github.com/susek555/BD2/car-dealer-api/internal/domains/user"

type Review struct {
	ID          int64      `json:"id" gorm:"primary_key"`
	Description string     `json:"description" gorm:"not null"`
	Rating      int        `json:"rating" gorm:"not null"`
	Reviewer    *user.User `gorm:"foreignKey:ReviewerID;references:ID"`
	Reviewed    *user.User `gorm:"foreignKey:ReviewedID;references:ID"`
}
