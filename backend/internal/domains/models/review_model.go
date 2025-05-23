package models

type Review struct {
	ID          uint   `json:"id" gorm:"primary_key;autoIncrement"`
	Description string `json:"description" gorm:"not null"`
	ReviewDate  string `json:"date" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Rating      uint   `json:"rating" gorm:"not null;check:rating BETWEEN 1 AND 5"`
	ReviewerID  uint   `json:"reviewer_id" gorm:"not null;check:reviewer_id <> reviewee_id;uniqueIndex:ux_review_pair"`
	Reviewer    *User  `gorm:"foreignKey:ReviewerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RevieweeId  uint   `json:"reviewee_id" gorm:"not null;uniqueIndex:ux_review_pair"`
	Reviewee    *User  `gorm:"foreignKey:RevieweeId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
