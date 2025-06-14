package sale_offer

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

var (
	ErrMissingFields                = errors.New("some fields are missing - ensure that all required fields are present")
	ErrInvalidProductionYear        = errors.New("invalid production year, it cannot be greater than current year")
	ErrInvalidRegistrationDate      = errors.New("invalid registration date, it cannot be in the future")
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
	ErrInvalidManufacturerModelPair = errors.New("manufacturer and model don't match")
	ErrOfferNotOwned                = errors.New("provided offer does not belong to logged in user")
	ErrOfferNotReadyToPublish       = errors.New("offer is not ready to be published - make sure that it have at least 3 images attached")
	ErrOfferOwnedByUser             = errors.New("offer is owned by user - cannot buy your own offer")
	ErrOfferAlreadySold             = errors.New("offer is already sold - cannot buy it again")
	ErrOfferNotPublished            = errors.New("offer is not published - cannot buy it")
	ErrOfferIsAuction               = errors.New("offer is an auction - cannot buy it directly, use bids instead")
	ErrOfferHasBids                 = errors.New("offer already has some bids - it cannot be updated/deleted")
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
	ErrOfferAlreadySold:             http.StatusConflict,
	ErrOfferNotPublished:            http.StatusBadRequest,
	ErrOfferIsAuction:               http.StatusBadRequest,
}
