package review

import "github.com/susek555/BD2/car-dealer-api/internal/domains/user"

type ReviewDTO struct {
	ID          uint    `json:"id"`
	Description string  `json:"description"`
	Rating      uint    `json:"rating"`
	Reviewer    UserDTO `json:"reviewer"`
	Reviewee    UserDTO `json:"reviewee"`
}

type UserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func (r *Review) MapToDTO() ReviewDTO {
	reviewDTO := ReviewDTO{}
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
