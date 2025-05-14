package sale_offer

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidColor         error = errors.New("invalid color")
	ErrInvalidFuelType      error = errors.New("invalid fuel type")
	ErrInvalidTransmission  error = errors.New("invalid transmission")
	ErrInvalidDrive         error = errors.New("invalid drive")
	ErrInvalidSaleOfferType error = errors.New("invalid sale offer type")
	ErrInvalidRange         error = errors.New("the min value should be lower than max")
	ErrInvalidDateFromat    error = errors.New("invalid date format, should be YYYY-MM-DD")
	ErrInvalidOrderKey      error = errors.New("invalid order-key")
	ErrInvalidManufacturer  error = errors.New("invalid manufacturer")
)

var ErrorMap = map[error]int{
	ErrInvalidColor:    http.StatusBadRequest,
	ErrInvalidDrive:    http.StatusBadRequest,
	ErrInvalidFuelType: http.StatusBadRequest,
	ErrInvalidDrive:    http.StatusBadRequest,
}
