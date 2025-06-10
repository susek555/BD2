package auth

type RegisterResponse struct {
	Errors map[string][]string `json:"errors,omitempty"`
}
