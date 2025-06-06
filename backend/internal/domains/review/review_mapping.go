package review

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/models"
)

func (ri *CreateReviewDTO) MapToObject(reviewerID uint) (*models.Review, error) {
	if !valIDateRating(ri.Rating) {
		return nil, ErrInvalidRating
	}
	return &models.Review{
		Description: ri.Description,
		Rating:      ri.Rating,
		ReviewerID:  reviewerID,
		RevieweeID:  ri.RevieweeID,
		ReviewDate:  time.Now(),
	}, nil
}

func (ur *UpdateReviewDTO) MapToObject(reviewerID, revieweeID uint) (*models.Review, error) {
	if !valIDateRating(ur.Rating) {
		return nil, ErrInvalidRating
	}
	return &models.Review{
		ID:          ur.ID,
		Description: ur.Description,
		Rating:      ur.Rating,
		ReviewerID:  reviewerID,
		RevieweeID:  revieweeID,
		ReviewDate:  time.Now(),
	}, nil
}

func MapToDTO(r *models.Review) *RetrieveReviewDTO {
	reviewDTO := RetrieveReviewDTO{}
	reviewDTO.ID = r.ID
	reviewDTO.Description = r.Description
	reviewDTO.Rating = r.Rating
	reviewee := MapToUserDTO(r.Reviewee)
	reviewDTO.Reviewee = reviewee
	reviewer := MapToUserDTO(r.Reviewer)
	reviewDTO.Reviewer = reviewer
	reviewDTO.ReviewDate = r.ReviewDate.Format(time.RFC3339)
	return &reviewDTO
}

func MapToUserDTO(u *models.User) UserDTO {
	return UserDTO{
		ID:       u.ID,
		Username: u.Username,
	}
}

func valIDateRating(rating uint) bool {
	return rating >= 1 && rating <= 5
}
