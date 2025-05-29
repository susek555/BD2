package image

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

var (
	ErrTooManyImages = errors.New("you can only attach 10 photos to one offer")
	ErrOfferNotOwned = errors.New("offer does not belong to logged-in user")
	ErrZeroImages    = errors.New("offer has no images - there is nothing to delete")
)

var ErrorMap = map[error]int{
	ErrTooManyImages:       http.StatusBadRequest,
	ErrZeroImages:          http.StatusBadRequest,
	ErrOfferNotOwned:       http.StatusForbidden,
	gorm.ErrRecordNotFound: http.StatusNotFound,
}
