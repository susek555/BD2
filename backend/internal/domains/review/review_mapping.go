package review

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"time"
)

func (ri *CreateReviewDTO) MapToObject(reviewerId uint) (*models.Review, error) {
	if !validateRating(ri.Rating) {
		return nil, ErrInvalidRating
	}
	return &models.Review{
		Description: ri.Description,
		Rating:      ri.Rating,
		ReviewerID:  reviewerId,
		RevieweeId:  ri.RevieweeId,
		ReviewDate:  time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func (ur *UpdateReviewDTO) MapToObject(reviewerId, revieweeId uint) (*models.Review, error) {
	if !validateRating(ur.Rating) {
		return nil, ErrInvalidRating
	}
	return &models.Review{
		ID:          ur.ID,
		Description: ur.Description,
		Rating:      ur.Rating,
		ReviewerID:  reviewerId,
		RevieweeId:  revieweeId,
		ReviewDate:  time.Now().Format("2006-01-02 15:04:05"),
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
	reviewDTO.ReviewDate = r.ReviewDate
	return &reviewDTO
}

func MapToUserDTO(u *models.User) UserDTO {
	return UserDTO{
		ID:       u.ID,
		Username: u.Username,
	}
}

func validateRating(rating uint) bool {
	return rating >= 1 && rating <= 5
}
