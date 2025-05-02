package user

import "errors"

var ErrInvalidSelector error = errors.New("User selector has to be P (person) or C (company)")

func (dto *UserDTO) MapToUser() (User, error) {
	switch dto.Selector {
	case "P":
		return User{
				ID:       dto.ID,
				Username: dto.Username,
				Password: dto.Password,
				Email:    dto.Email,
				Person:   &Person{Name: dto.PersonName, Surname: dto.PersonSurname},
			},
			nil
	case "C":
		return User{
				ID:       dto.ID,
				Username: dto.Username,
				Password: dto.Password,
				Email:    dto.Email,
				Company:  &Company{Name: dto.CompanyName, NIP: dto.CompanyNIP},
			},
			nil
	default:
		return User{}, ErrInvalidSelector
	}
}

func (user *User) MapToDTO() (UserDTO, error) {
	switch user.Selector {
	case "P":
		return UserDTO{
				ID:            user.ID,
				Username:      user.Username,
				Password:      user.Password,
				Email:         user.Email,
				PersonName:    user.Person.Name,
				PersonSurname: user.Person.Surname,
			},
			nil
	case "C":
		return UserDTO{
				ID:          user.ID,
				Username:    user.Username,
				Password:    user.Password,
				Email:       user.Email,
				CompanyName: user.Company.Name,
				CompanyNIP:  user.Company.NIP,
			},
			nil
	default:
		return UserDTO{}, ErrInvalidSelector
	}

}
