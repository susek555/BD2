package views

import "gorm.io/gorm"

type UserOfferRepositoryInterface interface {
	GetUserInteractionsByOfferID(offerID uint) ([]UserOfferRecord, error)
	GetUserInteractionsByUserID(userID uint) ([]UserOfferRecord, error)
}

type UserOfferRepository struct {
	DB *gorm.DB
}

func NewUserOfferRepository(db *gorm.DB) UserOfferRepositoryInterface {
	return &UserOfferRepository{
		DB: db,
	}
}

func (r *UserOfferRepository) GetUserInteractionsByOfferID(offerID uint) ([]UserOfferRecord, error) {
	var records []UserOfferRecord
	err := r.DB.
		Table("user_offer_interactions").
		Where("offer_id = ?", offerID).
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *UserOfferRepository) GetUserInteractionsByUserID(userID uint) ([]UserOfferRecord, error) {
	var records []UserOfferRecord
	err := r.DB.
		Table("user_offer_interactions").
		Where("user_id = ?", userID).
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}
