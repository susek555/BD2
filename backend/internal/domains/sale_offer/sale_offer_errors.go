package sale_offer

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

var (
	ErrMissingFields                = errors.New("some fields are missing - ensure that all required fields are present")
	ErrInvalidColor                 = errors.New("invalid color")
	ErrInvalidFuelType              = errors.New("invalid fuel type")
	ErrInvalidTransmission          = errors.New("invalid transmission")
	ErrInvalidDrive                 = errors.New("invalid drive")
	ErrInvalidSaleOfferType         = errors.New("invalid sale offer type")
	ErrInvalidMargin                = errors.New("invalid margin")
	ErrInvalidRange                 = errors.New("the min value should be lower than max")
	ErrInvalidDateFormat            = errors.New("invalid date format, should be YYYY-MM-DD")
	ErrInvalidOrderKey              = errors.New("invalid order-key")
	ErrInvalidManufacturer          = errors.New("invalid manufacturer")
	ErrInvalidManufacturerModelPair = errors.New("manufcaturer and model don't match")
	ErrOfferNotOwned                = errors.New("provided offer does not belong to logged in user")
	ErrOfferNotReadyToPublish       = errors.New("offer is not ready to be published - make sure that it have at least 3 images attached")
	ErrOfferOwnedByUser             = errors.New("offer is owned by user - cannot buy your own offer")
)

var ErrorMap = map[error]int{
	ErrMissingFields:                http.StatusBadRequest,
	ErrInvalidColor:                 http.StatusBadRequest,
	ErrInvalidFuelType:              http.StatusBadRequest,
	ErrInvalidTransmission:          http.StatusBadRequest,
	ErrInvalidDrive:                 http.StatusBadRequest,
	ErrInvalidSaleOfferType:         http.StatusBadRequest,
	ErrInvalidMargin:                http.StatusBadRequest,
	ErrInvalidRange:                 http.StatusBadRequest,
	ErrInvalidDateFormat:            http.StatusBadRequest,
	ErrInvalidOrderKey:              http.StatusBadRequest,
	ErrInvalidManufacturer:          http.StatusBadRequest,
	ErrInvalidManufacturerModelPair: http.StatusBadRequest,
	ErrOfferNotOwned:                http.StatusForbidden,
	gorm.ErrRecordNotFound:          http.StatusNotFound,
	ErrOfferOwnedByUser:             http.StatusForbidden,
}
