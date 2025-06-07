package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/views"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

//go:generate mockery --name=SaleOfferRepositoryInterface --output=../../test/mocks --case=snake --with-expecter
type SaleOfferRepositoryInterface interface {
	Create(offer *models.SaleOffer) error
	Update(offer *models.SaleOffer) error
	UpdateStatus(offer *models.SaleOffer, status enums.Status) error
	GetByID(id uint) (*models.SaleOffer, error)
	GetViewByID(id uint) (*views.SaleOfferView, error)
	GetFiltered(filter OfferFilterInterface, pagRequest *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error)
	GetAllActiveAuctions() ([]views.SaleOfferView, error)
	Delete(id uint) error
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
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(offer).Error
}

func (r *SaleOfferRepository) UpdateStatus(offer *models.SaleOffer, status enums.Status) error {
	offer.Status = status
	return r.Update(offer)
}

func (r *SaleOfferRepository) GetByID(id uint) (*models.SaleOffer, error) {
	var offer models.SaleOffer
	err := r.DB.
		Preload("Car.Model.Manufacturer").
		Preload("Auction").
		First(&offer, id).Error
	return &offer, err
}

func (r *SaleOfferRepository) GetViewByID(id uint) (*views.SaleOfferView, error) {
	var offerView views.SaleOfferView
	err := r.DB.Table("sale_offer_view").First(&offerView, id).Error
	return &offerView, err
}

func (r *SaleOfferRepository) GetByUserID(id uint, pagRequest *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error) {
	saleOffers, paginationResponse, err := pagination.PaginateResults[views.SaleOfferView](pagRequest, r.DB.Table("sale_offer_view").Where("user_id = ?", id))
	if err != nil {
		return nil, nil, err
	}
	return saleOffers, paginationResponse, nil
}

func (r *SaleOfferRepository) GetFiltered(filter OfferFilterInterface, pagRequest *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error) {
	query := r.DB.Table("sale_offer_view")
	query, err := filter.ApplyOfferFilters(query)
	if err != nil {
		return nil, nil, err
	}
	saleOffers, paginationResponse, err := pagination.PaginateResults[views.SaleOfferView](pagRequest, query)
	if err != nil {
		return nil, nil, err
	}
	return saleOffers, paginationResponse, nil
}

func (r *SaleOfferRepository) GetAllActiveAuctions() ([]views.SaleOfferView, error) {
	var auctions []views.SaleOfferView
	err := r.DB.Table("sale_offer_view").Where("is_auction IS TRUE").Find(&auctions).Error
	if err != nil {
		return nil, err
	}
	return auctions, nil
}

func (r *SaleOfferRepository) Delete(id uint) error {
	return r.DB.Delete(&models.SaleOffer{}, id).Error
}
