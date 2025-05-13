package review

import "github.com/susek555/BD2/car-dealer-api/internal/domains/user"

func (ri *CreateReviewDTO) MapToObject(reviewerId uint) Review {
	return Review{
		Description: ri.Description,
		Rating:      ri.Rating,
		ReviewerID:  reviewerId,
		RevieweeId:  ri.RevieweeId,
	}
}

func (ur *UpdateReviewDTO) MapToObject(reviewerId uint) Review {
	return Review{
		ID:          ur.ID,
		Description: ur.Description,
		Rating:      ur.Rating,
		ReviewerID:  reviewerId,
	}
}

func (r *Review) MapToDTO() RetrieveReviewDTO {
	reviewDTO := RetrieveReviewDTO{}
	reviewDTO.ID = r.ID
	reviewDTO.Description = r.Description
	reviewDTO.Rating = r.Rating
	reviewee := MapToUserDTO(r.Reviewee)
	reviewDTO.Reviewee = reviewee
	reviewer := MapToUserDTO(r.Reviewer)
	reviewDTO.Reviewer = reviewer
	return reviewDTO
}

func MapToUserDTO(u *user.User) UserDTO {
	return UserDTO{
		ID:       u.ID,
		Username: u.Username,
	}
}
