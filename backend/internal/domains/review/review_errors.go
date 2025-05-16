package review

import "errors"

var (
	ErrNoReviewsFound = errors.New("no reviews found for this reviewee")
	ErrNotReviewer    = errors.New("you are not the reviewer of this review")
)
