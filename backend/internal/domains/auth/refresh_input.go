package auth

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
