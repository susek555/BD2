package review

import "github.com/susek555/BD2/car-dealer-api/pkg/pagination"

type RetrieveReviewDTO struct {
	ID          uint    `json:"id"`
	Description string  `json:"description"`
	Rating      uint    `json:"rating"`
	Reviewer    UserDTO `json:"reviewer"`
	Reviewee    UserDTO `json:"reviewee"`
	ReviewDate  string  `json:"review_date"`
}

type CreateReviewDTO struct {
	Description string `json:"description"`
	Rating      uint   `json:"rating"`
	RevieweeId  uint   `json:"reviewee_id"`
}

type UpdateReviewDTO struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	Rating      uint   `json:"rating"`
}

type UserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type RetrieveReviewsWithPagination struct {
	Reviews            []RetrieveReviewDTO            `json:"reviews"`
	PaginationResponse *pagination.PaginationResponse `json:"pagination"`
}
