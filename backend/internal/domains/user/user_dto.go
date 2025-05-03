package user

type CreateUserDTO struct {
	Username      string  `json:"username"`
	Password      string  `json:"password"`
	Email         string  `json:"email"`
	Selector      string  `json:"selector"`
	CompanyName   *string `json:"company_name"`
	CompanyNIP    *string `json:"company_nip"`
	PersonName    *string `json:"person_name"`
	PersonSurname *string `json:"person_surname"`
}

type RetrieveUserDTO struct {
	Username      string  `json:"username"`
	Email         string  `json:"email"`
	CompanyName   *string `json:"company_name,omitempty"`
	CompanyNIP    *string `json:"company_nip,omitempty"`
	PersonName    *string `json:"person_name,omitempty"`
	PersonSurname *string `json:"person_surname,omitempty"`
}

type UpdateUserDTO struct {
	ID       uint    `json:"id"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
}
