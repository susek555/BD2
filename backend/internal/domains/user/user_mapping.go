package user

import (
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
)

func (dto *CreateUserDTO) MapToUser() (*User, error) {
	hashed, err := passwords.Hash(dto.Password)
	if err != nil {
		return nil, ErrHashPassword
	}
	if err := dto.validate(); err != nil {
		return nil, err
	}
	switch dto.Selector {
	case "P":
		if err := dto.validateP(); err != nil {
			return nil, err
		}
		return &User{
				Username: dto.Username,
				Password: hashed,
				Email:    dto.Email,
				Selector: dto.Selector,
				Person:   &Person{Name: *dto.PersonName, Surname: *dto.PersonSurname},
			},
			nil
	case "C":
		if err := dto.validateC(); err != nil {
			return nil, err
		}
		return &User{
				Username: dto.Username,
				Password: hashed,
				Email:    dto.Email,
				Selector: dto.Selector,
				Company:  &Company{Name: *dto.CompanyName, NIP: *dto.CompanyNIP},
			},
			nil
	default:
		return nil, ErrInvalidSelector
	}
}
func (dto *CreateUserDTO) validate() error {
	if dto.Username == "" || dto.Password == "" || dto.Email == "" || dto.Selector == "" {
		return ErrCreateUser
	}
	return nil
}
func (dto *CreateUserDTO) validateP() error {
	if dto.PersonName == nil || dto.PersonSurname == nil {
		return ErrCreatePerson
	}
	return nil
}

func (dto *CreateUserDTO) validateC() error {
	if dto.CompanyName == nil || dto.CompanyNIP == nil {
		return ErrCreateCompany
	}
	return nil
}

func (user *User) MapToDTO() *RetrieveUserDTO {
	switch user.Selector {
	case "P":
		return &RetrieveUserDTO{
			ID:            user.ID,
			Username:      user.Username,
			Email:         user.Email,
			PersonName:    &user.Person.Name,
			PersonSurname: &user.Person.Surname,
		}
	case "C":
		return &RetrieveUserDTO{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			CompanyName: &user.Company.Name,
			CompanyNIP:  &user.Company.NIP,
		}
	}
	return nil
}

func (dto *UpdateUserDTO) UpdateUserFromDTO(user *User) (*User, error) {
	if dto.Email != nil {
		user.Email = *dto.Email
	}
	if dto.Password != nil {
		newPassword, err := passwords.Hash(*dto.Password)
		if err != nil {
			return nil, ErrHashPassword
		}
		user.Password = newPassword
	}
	if dto.Username != nil {
		user.Username = *dto.Username
	}
	return user, nil
}
