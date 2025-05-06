package auth

type LoginResponse struct {
	RefreshToken  string              `json:"refresh_token,omitempty"`
	AccessToken   string              `json:"access_token,omitempty"`
	Selector      string              `json:"selector,omitempty"`
	Username      string              `json:"username,omitempty"`
	Email         string              `json:"email,omitempty"`
	PersonName    string              `json:"person_name,omitempty"`
	PersonSurname string              `json:"person_surname,omitempty"`
	CompanyName   string              `json:"company_name,omitempty"`
	CompanyNip    string              `json:"company_nip,omitempty"`
	Errors        map[string][]string `json:"errors,omitempty"`
}
