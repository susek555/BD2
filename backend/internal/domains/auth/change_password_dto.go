package auth

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
type ChangePasswordResponse struct {
	Errors map[string][]string `json:"errors"`
}
