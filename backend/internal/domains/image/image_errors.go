package image

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

var (
	ErrInvalidOfferID = errors.New("offer does not belong to logged-in user")
	ErrTooManyImages  = errors.New("you can only attach 10 photos to one offer")
)

var ErrorMap = map[error]int{
	ErrInvalidOfferID:      http.StatusForbidden,
	gorm.ErrRecordNotFound: http.StatusNotFound,
}
