package auth

type LoginInput struct {
	Login    string `json:"login" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
