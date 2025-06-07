package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"gorm.io/gorm"
)

type OfferFilterInterface interface {
	ApplyOfferFilters(*gorm.DB) (*gorm.DB, error)
	GetBase() *BaseOfferFilter
}

type PublishedOffersOnlyFilter struct {
	BaseOfferFilter
}

func (f *PublishedOffersOnlyFilter) ApplyOfferFilters(query *gorm.DB) (*gorm.DB, error) {
	query, err := f.BaseOfferFilter.ApplyOfferFilters(query)
	if err != nil {
		return nil, err
	}
	return applyPublishedOffersOnly(query, f.UserID), nil
}

func (f *PublishedOffersOnlyFilter) GetBase() *BaseOfferFilter {
	return &f.BaseOfferFilter
}

func applyPublishedOffersOnly(query *gorm.DB, userID *uint) *gorm.DB {
	query = query.Where("sale_offer_view.status = ?", enums.PUBLISHED)
	if userID != nil {
		query = query.Where("sale_offer_view.user_id = ?", *userID)
	}
	return query
}

type LikedOffersOnlyFilter struct {
	BaseOfferFilter
}

func (f *LikedOffersOnlyFilter) ApplyOfferFilters(query *gorm.DB) (*gorm.DB, error) {
	query, err := f.BaseOfferFilter.ApplyOfferFilters(query)
	if err != nil {
		return nil, err
	}
	return applyLikedOffersOnlyFilter(query, *f.UserID), nil
}

func (f *LikedOffersOnlyFilter) GetBase() *BaseOfferFilter {
	return &f.BaseOfferFilter
}

func applyLikedOffersOnlyFilter(query *gorm.DB, userID uint) *gorm.DB {
	return query.
		Joins("JOIN liked_offers ON liked_offers.offer_id = sale_offer_view.id").
		Where("liked_offers.user_id = ?", userID)
}

type UsersOffersOnlyFilter struct {
	BaseOfferFilter
}

func (f *UsersOffersOnlyFilter) ApplyOfferFilters(query *gorm.DB) (*gorm.DB, error) {
	query, err := f.BaseOfferFilter.ApplyOfferFilters(query)
	if err != nil {
		return nil, err
	}
	return applyUsersOffersOnly(query, *f.UserID), nil
}

func (f *UsersOffersOnlyFilter) GetBase() *BaseOfferFilter {
	return &f.BaseOfferFilter
}

func applyUsersOffersOnly(query *gorm.DB, userID uint) *gorm.DB {
	return query.Where("sale_offer_view.user_id = ?", userID)
}
