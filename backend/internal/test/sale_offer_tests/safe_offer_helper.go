package sale_offer_tests

import (
	"reflect"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ------
// Setup
// ------

func setupDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(
		&manufacturer.Manufacturer{},
		&model.Model{},
		&sale_offer.Car{},
		&sale_offer.SaleOffer{},
		&sale_offer.Auction{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getRepositoryWithSaleOffers(db *gorm.DB, offers []sale_offer.SaleOffer) sale_offer.SaleOfferRepositoryInterface {
	repo := sale_offer.NewSaleOfferRepository(db)
	for _, offer := range offers {
		repo.Create(&offer)
	}
	return repo
}

type OfferOption func(*sale_offer.SaleOffer, *sale_offer.Car)

// Simulates interaction with manufacturer service, which should return all possible manufacturers
var manufacturers []string = []string{"Audi", "BMW", "Opel", "Toyota", "Skoda"}

func CreateOffer(id uint, options ...OfferOption) *sale_offer.SaleOffer {
	car := &sale_offer.Car{
		OfferID:            id,
		Vin:                "vin",
		ProductionYear:     2025,
		Mileage:            1000,
		NumberOfDoors:      4,
		NumberOfSeats:      5,
		EnginePower:        100,
		EngineCapacity:     2000,
		RegistrationNumber: "default",
		RegistrationDate:   time.Now(),
		Color:              car_params.BLACK,
		FuelType:           car_params.PETROL,
		Transmission:       car_params.MANUAL,
		NumberOfGears:      6,
		Drive:              car_params.FWD,
		ModelID:            id,
		Model: model.Model{
			ID:             id,
			Name:           "model",
			ManufacturerID: id,
			Manufacturer: manufacturer.Manufacturer{
				ID:   id,
				Name: manufacturers[id-1],
			},
		},
	}
	offer := &sale_offer.SaleOffer{
		ID:          id,
		Description: "offer",
		Price:       1000,
		Margin:      15,
		DateOfIssue: time.Now(),
		Car:         car,
	}
	for _, option := range options {
		option(offer, car)
	}
	return offer
}

func WithCarField(fieldName string, fieldValue interface{}) OfferOption {
	return func(_ *sale_offer.SaleOffer, car *sale_offer.Car) {
		v := reflect.ValueOf(car).Elem()
		field := v.FieldByName(fieldName)
		field.Set(reflect.ValueOf(fieldValue))
	}
}

func WithOfferField(fieldName string, fieldValue interface{}) OfferOption {
	return func(offer *sale_offer.SaleOffer, _ *sale_offer.Car) {
		v := reflect.ValueOf(offer).Elem()
		field := v.FieldByName(fieldName)
		field.Set(reflect.ValueOf(fieldValue))
	}
}

func WithAuction(dateEnd time.Time, buyNowPrice uint) OfferOption {
	return func(offer *sale_offer.SaleOffer, _ *sale_offer.Car) {
		offer.Auction = &sale_offer.Auction{
			OfferID:     offer.ID,
			DateEnd:     dateEnd,
			BuyNowPrice: buyNowPrice,
		}
	}
}

func GetDefaultPaginationRequest() *pagination.PaginationRequest {
	return &pagination.PaginationRequest{Page: 1, PageSize: 8}
}
