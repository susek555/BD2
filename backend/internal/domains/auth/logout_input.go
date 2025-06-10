package auth

type LogoutInput struct {
	RefreshToken string `json:"refresh_token"`
	AllDevices   bool   `json:"all_devices"`
}
