package user

import (
	"errors"

	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
)

var (
	ErrInvalidSelector error = errors.New("user selector has to be P (person) or C (company)")
	ErrHashPassword    error = errors.New("error occured while hashing password")
)

func (dto *CreateUserDTO) MapToUser() (User, error) {
	switch dto.Selector {
	case "P":
		return User{
				Username: dto.Username,
				Password: dto.Password,
				Email:    dto.Email,
				Selector: dto.Selector,
				Person:   &Person{Name: *dto.PersonName, Surname: *dto.PersonSurname},
			},
			nil
	case "C":
		return User{
				Username: dto.Username,
				Password: dto.Password,
				Email:    dto.Email,
				Selector: dto.Selector,
				Company:  &Company{Name: *dto.CompanyName, NIP: *dto.CompanyNIP},
			},
			nil
	default:
		return User{}, ErrInvalidSelector
	}
}

func (user *User) MapToDTO() (RetrieveUserDTO, error) {
	switch user.Selector {
	case "P":
		return RetrieveUserDTO{
				Username:      user.Username,
				Email:         user.Email,
				PersonName:    &user.Person.Name,
				PersonSurname: &user.Person.Surname,
			},
			nil
	case "C":
		return RetrieveUserDTO{
				Username:    user.Username,
				Email:       user.Email,
				CompanyName: &user.Company.Name,
				CompanyNIP:  &user.Company.NIP,
			},
			nil
	default:
		return RetrieveUserDTO{}, ErrInvalidSelector
	}

}

func (dto *UpdateUserDTO) UpdateUserFromDTO(user *User) (*User, error) {
	if dto.Email != nil {
		user.Email = *dto.Email
	}
	if dto.Password != nil {
		newPassword, err := passwords.Hash(*dto.Password)
		if err != nil {
			return &User{}, ErrHashPassword
		}
		user.Password = newPassword
	}
	if dto.Username != nil {
		user.Username = *dto.Username
	}
	return user, nil
}
