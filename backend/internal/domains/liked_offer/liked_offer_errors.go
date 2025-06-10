package liked_offer

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

var (
	ErrLikeAlreadyLikedOffer = errors.New("offer already liked")
	ErrDislikeNotLikedOffer  = errors.New("offer not liked before cannot be disliked")
	ErrLikeOwnOffer          = errors.New("your own offer cannot be liked")
)

var ErrorMap = map[error]int{
	ErrLikeAlreadyLikedOffer: http.StatusBadRequest,
	ErrDislikeNotLikedOffer:  http.StatusBadRequest,
	ErrLikeOwnOffer:          http.StatusBadRequest,
	gorm.ErrRecordNotFound:   http.StatusNotFound,
}
