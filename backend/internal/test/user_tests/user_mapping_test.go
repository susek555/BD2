//go:build unit
// +build unit

package user_tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	u "github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
)

func createPerson() *u.User {
	user := u.User{
		ID:       1,
		Username: "john",
		Email:    "john@example.com",
		Password: "hashed_password",
		Selector: "P",
		Person:   &u.Person{Name: "john person", Surname: "doe person"},
	}
	return &user
}

func createCompany() *u.User {
	user := u.User{
		ID:       1,
		Username: "john",
		Email:    "john@example.com",
		Password: "hashed_password",
		Selector: "C",
		Company:  &u.Company{Name: "john company", NIP: "1234567890"},
	}
	return &user
}

func TestMapToUser_EmptyDTO(t *testing.T) {
	dto := u.CreateUserDTO{}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, u.ErrCreateUser)
}

func TestMapToUser_MissingUsername(t *testing.T) {
	dto := u.CreateUserDTO{Password: "123", Email: "john@example.com", Selector: "P"}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, u.ErrCreateUser)
}

func TestMapToUser_MissingPassword(t *testing.T) {
	dto := u.CreateUserDTO{Username: "john", Email: "john@exampl.com", Selector: "P"}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, u.ErrCreateUser)
}

func TestMapToUser_MissingEmail(t *testing.T) {
	dto := u.CreateUserDTO{Username: "john", Password: "123", Selector: "P"}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, u.ErrCreateUser)
}

func TestMapToUser_MissingSelector(t *testing.T) {
	dto := u.CreateUserDTO{Username: "john", Password: "123", Email: "john@example.com"}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, u.ErrCreateUser)
}

func TestMapToUser_InvalidSelector(t *testing.T) {
	dto := u.CreateUserDTO{Username: "john", Password: "123", Email: "john@example.com", Selector: "X"}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, u.ErrInvalidSelector)
}
func TestMapToUser_EmptyPersonName(t *testing.T) {
	surname := "doe"
	dto := u.CreateUserDTO{
		Username:      "john",
		Password:      "123",
		Email:         "john@example.com",
		Selector:      "P",
		PersonName:    nil,
		PersonSurname: &surname,
	}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, u.ErrCreatePerson)
}

func TestMapToUser_EmptyPersonSurname(t *testing.T) {
	name := "john person"
	dto := u.CreateUserDTO{
		Username:      "john",
		Password:      "123",
		Email:         "john@example.com",
		Selector:      "P",
		PersonName:    &name,
		PersonSurname: nil,
	}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, u.ErrCreatePerson)
}

func TestMapToUser_PersonWithCompanyFields(t *testing.T) {
	name := "john company"
	nip := "1234567890"
	dto := u.CreateUserDTO{
		Username:    "john",
		Password:    "123",
		Email:       "john@example.com",
		Selector:    "P",
		CompanyName: &name,
		CompanyNIP:  &nip,
	}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, u.ErrCreatePerson)
}

func TestMapToUser_EmptyCompanyName(t *testing.T) {
	nip := "1234567890"
	dto := u.CreateUserDTO{
		Username:    "john",
		Password:    "123",
		Email:       "john@example.com",
		Selector:    "C",
		CompanyName: nil,
		CompanyNIP:  &nip,
	}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, u.ErrCreateCompany)
}

func TestMapToUser_EmptyCompanyNIP(t *testing.T) {
	name := "john Company"
	dto := u.CreateUserDTO{
		Username:    "john",
		Password:    "123",
		Email:       "john@example.com",
		Selector:    "C",
		CompanyName: &name,
		CompanyNIP:  nil,
	}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, u.ErrCreateCompany)
}

func TestMapToUser_CompanyWithPersonFields(t *testing.T) {
	name := "john person"
	surname := "doe person"
	dto := u.CreateUserDTO{
		Username:      "john",
		Password:      "123",
		Email:         "john@example.com",
		Selector:      "C",
		PersonName:    &name,
		PersonSurname: &surname,
	}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, u.ErrCreateCompany)
}

func TestMapToUser_ValidPerson(t *testing.T) {
	name := "john person"
	surname := "doe person"
	dto := u.CreateUserDTO{
		Username:      "john",
		Password:      "123",
		Email:         "john@example.com",
		Selector:      "P",
		PersonName:    &name,
		PersonSurname: &surname,
	}
	user, err := dto.MapToUser()
	assert.NoError(t, err)
	assert.Equal(t, "john", user.Username)
	assert.Equal(t, passwords.Match("123", user.Password), true)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "P", user.Selector)
	assert.Equal(t, "john person", user.Person.Name)
	assert.Equal(t, "doe person", user.Person.Surname)
}

func TestMapToUser_ValidCompany(t *testing.T) {
	name := "john company"
	nip := "1234567890"
	dto := u.CreateUserDTO{
		Username:    "john",
		Password:    "123",
		Email:       "john@example.com",
		Selector:    "C",
		CompanyName: &name,
		CompanyNIP:  &nip,
	}
	user, err := dto.MapToUser()
	assert.NoError(t, err)
	assert.Equal(t, "john", user.Username)
	assert.Equal(t, passwords.Match("123", user.Password), true)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "C", user.Selector)
	assert.Equal(t, "john company", user.Company.Name)
	assert.Equal(t, "1234567890", user.Company.NIP)
}

func TestMapToUser_PersonWithAllFields(t *testing.T) {
	companyName := "john company"
	companyNIP := "1234567890"
	name := "john person"
	surname := "doe person"
	dto := u.CreateUserDTO{
		Username:      "john",
		Password:      "123",
		Email:         "john@example.com",
		Selector:      "P",
		PersonName:    &name,
		PersonSurname: &surname,
		CompanyName:   &companyName,
		CompanyNIP:    &companyNIP,
	}
	user, err := dto.MapToUser()
	assert.NoError(t, err)
	assert.Equal(t, "john", user.Username)
	assert.Equal(t, passwords.Match("123", user.Password), true)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "P", user.Selector)
	assert.Equal(t, "john person", user.Person.Name)
	assert.Equal(t, "doe person", user.Person.Surname)
	assert.Nil(t, user.Company)
}

func TestMapToUser_CompanyWithAllFields(t *testing.T) {
	companyName := "john company"
	companyNIP := "1234567890"
	name := "john person"
	surname := "doe person"
	dto := u.CreateUserDTO{
		Username:      "john",
		Password:      "123",
		Email:         "john@example.com",
		Selector:      "C",
		CompanyName:   &companyName,
		CompanyNIP:    &companyNIP,
		PersonName:    &name,
		PersonSurname: &surname,
	}
	user, err := dto.MapToUser()
	assert.NoError(t, err)
	assert.Equal(t, "john", user.Username)
	assert.Equal(t, passwords.Match("123", user.Password), true)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "C", user.Selector)
	assert.Equal(t, "john company", user.Company.Name)
	assert.Equal(t, "1234567890", user.Company.NIP)
	assert.Nil(t, user.Person)
}

func TestMapToDTO_Person(t *testing.T) {
	user := createPerson()
	dto, err := user.MapToDTO()
	assert.NoError(t, err)
	assert.Equal(t, user.Username, dto.Username)
	assert.Equal(t, user.Email, dto.Email)
	assert.Equal(t, user.Person.Name, *dto.PersonName)
	assert.Equal(t, user.Person.Surname, *dto.PersonSurname)
}
func TestMapToDTO_Company(t *testing.T) {
	user := createCompany()
	dto, err := user.MapToDTO()
	assert.NoError(t, err)
	assert.Equal(t, user.Username, dto.Username)
	assert.Equal(t, user.Email, dto.Email)
	assert.Equal(t, user.Company.Name, *dto.CompanyName)
	assert.Equal(t, user.Company.NIP, *dto.CompanyNIP)
}

func TestUpdateUserFromDTO_Email(t *testing.T) {
	email := "john_updated@example.com"
	dto := u.UpdateUserDTO{
		ID:    1,
		Email: &email,
	}
	user := createPerson()
	new_user, _ := dto.UpdateUserFromDTO(user)
	assert.Equal(t, email, new_user.Email)
	assert.Equal(t, user.Username, new_user.Username)
	assert.Equal(t, user.Password, new_user.Password)
	assert.Equal(t, user.Selector, new_user.Selector)
}

func TestUpdateUserFromDTO_Password(t *testing.T) {
	password := "new_password"
	dto := u.UpdateUserDTO{
		ID:       1,
		Password: &password,
	}
	user := createPerson()
	new_user, _ := dto.UpdateUserFromDTO(user)
	assert.Equal(t, user.Email, new_user.Email)
	assert.Equal(t, user.Username, new_user.Username)
	assert.Equal(t, passwords.Match(password, new_user.Password), true)
	assert.Equal(t, user.Selector, new_user.Selector)
}

func TestUpdateUserFromDTO_Username(t *testing.T) {
	username := "john_updated"
	dto := u.UpdateUserDTO{
		ID:       1,
		Username: &username,
	}
	user := createPerson()
	new_user, _ := dto.UpdateUserFromDTO(user)
	assert.Equal(t, user.Email, new_user.Email)
	assert.Equal(t, user.Password, new_user.Password)
	assert.Equal(t, username, new_user.Username)
	assert.Equal(t, user.Selector, new_user.Selector)
}

func TestUpdateUserFromDTO_Empty(t *testing.T) {
	dto := u.UpdateUserDTO{
		ID: 1,
	}
	user := createPerson()
	new_user, _ := dto.UpdateUserFromDTO(user)
	assert.Equal(t, user.Email, new_user.Email)
	assert.Equal(t, user.Password, new_user.Password)
	assert.Equal(t, user.Username, new_user.Username)
	assert.Equal(t, user.Selector, new_user.Selector)
}
