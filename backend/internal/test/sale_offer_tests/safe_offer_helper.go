package sale_offer_tests

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	u "github.com/susek555/BD2/car-dealer-api/internal/test/test_utils"
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

func createOffer(id uint) *sale_offer.SaleOffer {
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
	return offer
}

func withCarField(opt u.Option[sale_offer.Car]) u.Option[sale_offer.SaleOffer] {
	return func(offer *sale_offer.SaleOffer) {
		if offer.Car == nil {
			offer.Car = &sale_offer.Car{}
		}
		opt(offer.Car)
	}
}

func withAuctionField(opt u.Option[sale_offer.Auction]) u.Option[sale_offer.SaleOffer] {
	return func(offer *sale_offer.SaleOffer) {
		if offer.Auction == nil {
			offer.Auction = &sale_offer.Auction{}
		}
		opt(offer.Auction)
	}
}

func createAuctionSaleOffer(id uint) *sale_offer.SaleOffer {
	offer := *u.Build(createOffer(id),
		withAuctionField(u.WithField[sale_offer.Auction]("OfferID", id)),
		withAuctionField(u.WithField[sale_offer.Auction]("DateEnd", time.Now())),
		withAuctionField(u.WithField[sale_offer.Auction]("BuyNowPrice", uint(0))))
	return &offer
}
