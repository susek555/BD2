package sale_offer

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

var (
	ErrInvalidColor         error = errors.New("invalid color")
	ErrInvalidFuelType      error = errors.New("invalid fuel type")
	ErrInvalidTransmission  error = errors.New("invalid transmission")
	ErrInvalidDrive         error = errors.New("invalid drive")
	ErrInvalidSaleOfferType error = errors.New("invalid sale offer type")
	ErrInvalidMargin        error = errors.New("invalid margin")
	ErrInvalidRange         error = errors.New("the min value should be lower than max")
	ErrInvalidDateFromat    error = errors.New("invalid date format, should be YYYY-MM-DD")
	ErrInvalidOrderKey      error = errors.New("invalid order-key")
	ErrInvalidManufacturer  error = errors.New("invalid manufacturer")
	ErrAuthorization        error = errors.New("you have to be logged in to create an offer")
)

var ErrorMap = map[error]int{
	ErrInvalidColor:         http.StatusBadRequest,
	ErrInvalidFuelType:      http.StatusBadRequest,
	ErrInvalidTransmission:  http.StatusBadRequest,
	ErrInvalidDrive:         http.StatusBadRequest,
	ErrInvalidSaleOfferType: http.StatusBadRequest,
	ErrInvalidMargin:        http.StatusBadRequest,
	ErrInvalidRange:         http.StatusBadRequest,
	ErrInvalidDateFromat:    http.StatusBadRequest,
	ErrInvalidOrderKey:      http.StatusBadRequest,
	ErrInvalidManufacturer:  http.StatusBadRequest,
	ErrAuthorization:        http.StatusUnauthorized,
	gorm.ErrRecordNotFound:  http.StatusNotFound,
}
