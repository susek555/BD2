package review

type RetrieveReviewDTO struct {
	ID          uint    `json:"id"`
	Description string  `json:"description"`
	Rating      uint    `json:"rating"`
	Reviewer    UserDTO `json:"reviewer"`
	Reviewee    UserDTO `json:"reviewee"`
	CreatedAt   string  `json:"created_at"`
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
