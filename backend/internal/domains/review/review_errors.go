package review

import (
	"errors"
	"net/http"
)

var (
	ErrNoReviewsFound = errors.New("no reviews found for this reviewee")
	ErrNotReviewer    = errors.New("you are not the reviewer of this review")
	ErrNoReviewFound  = errors.New("no review found")
	ErrInvalidRating  = errors.New("invalid rating, must be between 1 and 5")
)

var ErrorMap = map[error]int{
	ErrNoReviewsFound: http.StatusNotFound,
	ErrNotReviewer:    http.StatusForbidden,
	ErrNoReviewFound:  http.StatusNotFound,
	ErrInvalidRating:  http.StatusBadRequest,
}
