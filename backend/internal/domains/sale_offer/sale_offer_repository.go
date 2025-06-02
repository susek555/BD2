package sale_offer

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

type SaleOfferRepositoryInterface interface {
	Create(offer *models.SaleOffer) error
	Update(offer *models.SaleOffer) error
	GetFiltered(filter *OfferFilter) ([]SaleOfferView, *pagination.PaginationResponse, error)
	GetByID(id uint) (*SaleOfferView, error)
	GetByUserID(id uint, pagination *pagination.PaginationRequest) ([]SaleOfferView, *pagination.PaginationResponse, error)
	GetAllActiveAuctions() ([]models.SaleOffer, error)
	GetAllActiveOffers() ([]models.SaleOffer, error)
	BuyOffer(offerID uint, buyerID uint) (*SaleOfferView, error)
}

type SaleOfferRepository struct {
	DB *gorm.DB
}

func NewSaleOfferRepository(db *gorm.DB) SaleOfferRepositoryInterface {
	return &SaleOfferRepository{DB: db}
}

func (r *SaleOfferRepository) Create(offer *models.SaleOffer) error {
	return r.DB.Create(offer).Error
}

func (r *SaleOfferRepository) Update(offer *models.SaleOffer) error {
	return r.DB.Save(offer).Error
}

func (r *SaleOfferRepository) GetFiltered(filter *OfferFilter) ([]SaleOfferView, *pagination.PaginationResponse, error) {
	query := r.DB.Table("car_offer_view").Where("status = ?", enums.PUBLISHED)
	query, err := filter.ApplyOfferFilters(query)
	if err != nil {
		return nil, nil, err
	}
	saleOffers, paginationResponse, err := pagination.PaginateResults[SaleOfferView](&filter.Pagination, query)
	if err != nil {
		return nil, nil, err
	}
	return saleOffers, paginationResponse, nil
}

func (r *SaleOfferRepository) GetByID(id uint) (*SaleOfferView, error) {
	var offer SaleOfferView
	err := r.DB.Table("car_offer_view").Find(&offer, id).Error
	return &offer, err
}

func (r *SaleOfferRepository) GetByUserID(id uint, pagRequest *pagination.PaginationRequest) ([]SaleOfferView, *pagination.PaginationResponse, error) {
	saleOffers, paginationResponse, err := pagination.PaginateResults[SaleOfferView](pagRequest, r.DB.Table("car_offer_view").Where("user_id = ?", id))
	if err != nil {
		return nil, nil, err
	}
	return saleOffers, paginationResponse, nil
}

func (r *SaleOfferRepository) GetAllActiveAuctions() ([]models.SaleOffer, error) {
	var auctions []models.SaleOffer
	err := r.DB.
		Preload("Auction").
		Joins("JOIN auctions ON auctions.offer_id = sale_offers.id").
		Where("auctions.offer_id IS NOT NULL").
		Where("auctions.date_end > NOW()").
		Find(&auctions).
		Error
	if err != nil {
		return nil, err
	}
	return auctions, nil
}

func (r *SaleOfferRepository) GetAllActiveOffers() ([]models.SaleOffer, error) {
	var offers []models.SaleOffer
	err := r.DB.
		Preload("Auction").
		Joins("LEFT JOIN auctions ON auctions.offer_id = sale_offers.id").
		Where("auctions.date_end > NOW() OR auctions.offer_id IS NULL").
		Find(&offers).
		Error
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func (r *SaleOfferRepository) BuyOffer(offerID uint, buyerID uint) (*SaleOfferView, error) {
	offer, err := r.GetByID(offerID)
	if err != nil {
		return nil, err
	}
	if offer.Status == enums.SOLD {
		return nil, ErrOfferAlreadySold
	}
	err = r.DB.Model(&models.SaleOffer{}).
		Where("id = ?", offerID).
		Update("status", enums.SOLD).
		Error
	if err != nil {
		return nil, err
	}
	err = r.DB.Model(&models.Purchase{}).
		Create(&models.Purchase{
			OfferID:    offerID,
			BuyerID:    buyerID,
			FinalPrice: offer.Price,
			IssueDate:  time.Now(),
		}).Error
	return offer, err
}
