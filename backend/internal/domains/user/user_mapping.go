package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
)

func (dto *CreateUserDTO) MapToUser() (*User, error) {
	hashed, err := passwords.Hash(dto.Password)
	if err != nil {
		return nil, ErrHashPassword
	}
	validator := validator.New()
	err = validator.Struct(dto)
	if err != nil {
		return nil, ErrCreateUser
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
	if err := dto.updateMainFields(user); err != nil {
		return nil, err
	}
	if err := dto.updatePersonFields(user); err != nil {
		return nil, err
	}
	if err := dto.updateCompanyFields(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (dto *UpdateUserDTO) updateMainFields(user *User) error {
	if dto.Email != nil {
		user.Email = *dto.Email
	}
	if dto.Password != nil {
		newPassword, err := passwords.Hash(*dto.Password)
		if err != nil {
			return ErrHashPassword
		}
		user.Password = newPassword
	}
	if dto.Username != nil {
		user.Username = *dto.Username
	}
	return nil
}

func (dto *UpdateUserDTO) updatePersonFields(user *User) error {
	if dto.PersonName == nil && dto.PersonSurname == nil {
		return nil
	}
	if user.Selector != "P" {
		return ErrUpdatePerson
	}
	if dto.PersonName != nil {
		user.Person.Name = *dto.PersonName
	}
	if dto.PersonSurname != nil {
		user.Person.Surname = *dto.PersonSurname
	}
	return nil
}

func (dto *UpdateUserDTO) updateCompanyFields(user *User) error {
	if dto.CompanyName == nil && dto.CompanyNIP == nil {
		return nil
	}
	if user.Selector != "C" {
		return ErrUpdateCompany
	}
	if dto.CompanyName != nil {
		user.Company.Name = *dto.CompanyName
	}
	if dto.CompanyNIP != nil {
		user.Company.NIP = *dto.CompanyNIP
	}
	return nil
}
