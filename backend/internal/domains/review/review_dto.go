package review

type RetrieveReviewDTO struct {
	ID          uint    `json:"id"`
	Description string  `json:"description"`
	Rating      uint    `json:"rating"`
	Reviewer    UserDTO `json:"reviewer"`
	Reviewee    UserDTO `json:"reviewee"`
}

type CreateReviewDTO struct {
	Description string `json:"description"`
	Rating      uint   `json:"rating"`
	RevieweeId  uint   `json:"reviewee_id"`
}

type UserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}
