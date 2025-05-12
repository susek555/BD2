package sale_offer_tests

import (
	"reflect"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
)

type OfferOption func(*sale_offer.SaleOffer, *car.Car)

func createOffer(id uint, options ...OfferOption) *sale_offer.SaleOffer {
	car := &car.Car{
		ID:                 id,
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
		Model: &model.Model{
			ID:             id,
			Name:           "model",
			ManufacturerID: id,
			Manufacturer: manufacturer.Manufacturer{
				ID:   id,
				Name: "manufacturer",
			},
		},
	}
	offer := &sale_offer.SaleOffer{
		ID:          id,
		Description: "offer",
		Price:       1000,
		Margin:      15,
		DateOfIssue: time.Now(),
		CarID:       1,
		Car:         car,
	}
	for _, option := range options {
		option(offer, car)
	}
	return offer
}

func withCarField(fieldName string, fieldValue interface{}) OfferOption {
	return func(_ *sale_offer.SaleOffer, car *car.Car) {
		v := reflect.ValueOf(car).Elem()
		field := v.FieldByName(fieldName)
		field.Set(reflect.ValueOf(fieldValue))
	}
}

func withOfferField(fieldName string, fieldValue interface{}) OfferOption {
	return func(offer *sale_offer.SaleOffer, _ *car.Car) {
		v := reflect.ValueOf(offer).Elem()
		field := v.FieldByName(fieldName)
		field.Set(reflect.ValueOf(fieldValue))
	}
}

func withAuction(dateEnd time.Time, buyNowPrice uint) OfferOption {
	return func(offer *sale_offer.SaleOffer, _ *car.Car) {
		offer.Auction = &sale_offer.Auction{
			OfferID:     offer.ID,
			DateEnd:     dateEnd,
			BuyNowPrice: buyNowPrice,
		}
	}
}
