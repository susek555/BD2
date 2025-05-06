package auth

type LoginInput struct {
	Login    string `json:"login" binding:"required,login"`
	Password string `json:"password" binding:"required"`
}
