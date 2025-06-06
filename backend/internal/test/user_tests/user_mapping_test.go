package user_tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
)

// --------------------------------
// Test map to user from create dto
// --------------------------------

func TestMapToUser_EmptyDTO(t *testing.T) {
	dto := user.CreateUserDTO{}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, user.ErrCreateUser)
}

func TestMapToUser_MissingUsername(t *testing.T) {
	dto := user.CreateUserDTO{Password: "123", Email: "john@example.com", Selector: "P"}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, user.ErrCreateUser)
}

func TestMapToUser_MissingPassword(t *testing.T) {
	dto := user.CreateUserDTO{Username: "john", Email: "john@exampl.com", Selector: "P"}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, user.ErrCreateUser)
}

func TestMapToUser_MissingEmail(t *testing.T) {
	dto := user.CreateUserDTO{Username: "john", Password: "123", Selector: "P"}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, user.ErrCreateUser)
}

func TestMapToUser_MissingSelector(t *testing.T) {
	dto := user.CreateUserDTO{Username: "john", Password: "123", Email: "john@example.com"}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, user.ErrCreateUser)
}

func TestMapToUser_InvalidSelector(t *testing.T) {
	dto := user.CreateUserDTO{Username: "john", Password: "123", Email: "john@example.com", Selector: "X"}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, user.ErrInvalidSelector)
}
func TestMapToUser_EmptyPersonName(t *testing.T) {
	surname := "doe"
	dto := user.CreateUserDTO{
		Username:      "john",
		Password:      "123",
		Email:         "john@example.com",
		Selector:      "P",
		PersonName:    nil,
		PersonSurname: &surname,
	}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, user.ErrCreatePerson)
}

func TestMapToUser_EmptyPersonSurname(t *testing.T) {
	name := "john person"
	dto := user.CreateUserDTO{
		Username:      "john",
		Password:      "123",
		Email:         "john@example.com",
		Selector:      "P",
		PersonName:    &name,
		PersonSurname: nil,
	}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, user.ErrCreatePerson)
}

func TestMapToUser_PersonWithCompanyFields(t *testing.T) {
	name := "john company"
	nip := "1234567890"
	dto := user.CreateUserDTO{
		Username:    "john",
		Password:    "123",
		Email:       "john@example.com",
		Selector:    "P",
		CompanyName: &name,
		CompanyNIP:  &nip,
	}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, user.ErrCreatePerson)
}

func TestMapToUser_EmptyCompanyName(t *testing.T) {
	nip := "1234567890"
	dto := user.CreateUserDTO{
		Username:    "john",
		Password:    "123",
		Email:       "john@example.com",
		Selector:    "C",
		CompanyName: nil,
		CompanyNIP:  &nip,
	}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, user.ErrCreateCompany)
}

func TestMapToUser_EmptyCompanyNIP(t *testing.T) {
	name := "john Company"
	dto := user.CreateUserDTO{
		Username:    "john",
		Password:    "123",
		Email:       "john@example.com",
		Selector:    "C",
		CompanyName: &name,
		CompanyNIP:  nil,
	}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, user.ErrCreateCompany)
}

func TestMapToUser_CompanyWithPersonFields(t *testing.T) {
	name := "john person"
	surname := "doe person"
	dto := user.CreateUserDTO{
		Username:      "john",
		Password:      "123",
		Email:         "john@example.com",
		Selector:      "C",
		PersonName:    &name,
		PersonSurname: &surname,
	}
	_, err := dto.MapToUser()
	assert.ErrorIs(t, err, user.ErrCreateCompany)
}

func TestMapToUser_ValidPerson(t *testing.T) {
	name := "john person"
	surname := "doe person"
	dto := user.CreateUserDTO{
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
	dto := user.CreateUserDTO{
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
	assert.Equal(t, "1234567890", user.Company.Nip)
}

func TestMapToUser_PersonWithAllFields(t *testing.T) {
	companyName := "john company"
	companyNIP := "1234567890"
	name := "john person"
	surname := "doe person"
	dto := user.CreateUserDTO{
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
	dto := user.CreateUserDTO{
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
	assert.Equal(t, "1234567890", user.Company.Nip)
	assert.Nil(t, user.Person)
}

// ----------------------------------
// Test map to retrieve dto from user
// ----------------------------------

func TestMapToDTO_Person(t *testing.T) {
	user_ := createPerson(1)
	dto := user.MapToDTO(user_)
	assert.Equal(t, user_.ID, dto.ID)
	assert.Equal(t, user_.Username, dto.Username)
	assert.Equal(t, user_.Email, dto.Email)
	assert.Equal(t, user_.Person.Name, *dto.PersonName)
	assert.Equal(t, user_.Person.Surname, *dto.PersonSurname)
}
func TestMapToDTO_Company(t *testing.T) {
	user_ := createCompany(1)
	dto := user.MapToDTO(user_)
	assert.Equal(t, user_.ID, dto.ID)
	assert.Equal(t, user_.Username, dto.Username)
	assert.Equal(t, user_.Email, dto.Email)
	assert.Equal(t, user_.Company.Name, *dto.CompanyName)
	assert.Equal(t, user_.Company.Nip, *dto.CompanyNIP)
}

// --------------------------------
// Test update user from update dto
// --------------------------------

func TestUpdateUserFromDTO_Email(t *testing.T) {
	email := "john_updated@example.com"
	dto := user.UpdateUserDTO{
		ID:    1,
		Email: &email,
	}
	user := createPerson(1)
	newUser, _ := dto.UpdateUserFromDTO(user)
	assert.Equal(t, email, newUser.Email)
	assert.Equal(t, user.Username, newUser.Username)
	assert.Equal(t, user.Password, newUser.Password)
	assert.Equal(t, user.Selector, newUser.Selector)
}

func TestUpdateUserFromDTO_Password(t *testing.T) {
	password := "new_password"
	dto := user.UpdateUserDTO{
		ID:       1,
		Password: &password,
	}
	user := createPerson(1)
	newUser, _ := dto.UpdateUserFromDTO(user)
	assert.Equal(t, user.Email, newUser.Email)
	assert.Equal(t, user.Username, newUser.Username)
	assert.Equal(t, passwords.Match(password, newUser.Password), true)
	assert.Equal(t, user.Selector, newUser.Selector)
}

func TestUpdateUserFromDTO_Username(t *testing.T) {
	username := "john_updated"
	dto := user.UpdateUserDTO{
		ID:       1,
		Username: &username,
	}
	user := createPerson(1)
	newUser, _ := dto.UpdateUserFromDTO(user)
	assert.Equal(t, user.Email, newUser.Email)
	assert.Equal(t, user.Password, newUser.Password)
	assert.Equal(t, username, newUser.Username)
	assert.Equal(t, user.Selector, newUser.Selector)
}

func TestUpdateUserFromDTO_CompanyNameAsCompany(t *testing.T) {
	companyName := "new_name"
	dto := user.UpdateUserDTO{
		ID:          1,
		CompanyName: &companyName,
	}
	user := createCompany(1)
	newUser, _ := dto.UpdateUserFromDTO(user)
	assert.Equal(t, user.Email, newUser.Email)
	assert.Equal(t, user.Password, newUser.Password)
	assert.Equal(t, user.Username, newUser.Username)
	assert.Equal(t, companyName, newUser.Company.Name)
}

func TestUpdateUserFromDTO_CompanyNameAsPerson(t *testing.T) {
	companyName := "new_name"
	dto := user.UpdateUserDTO{
		ID:          1,
		CompanyName: &companyName,
	}
	user_ := createPerson(1)
	_, err := dto.UpdateUserFromDTO(user_)
	assert.ErrorIs(t, err, user.ErrUpdateCompany)
}

func TestUpdateUserFromDTO_CompanyNIPAsCompany(t *testing.T) {
	companyNIP := "new_nip"
	dto := user.UpdateUserDTO{
		ID:         1,
		CompanyNIP: &companyNIP,
	}
	user := createCompany(1)
	newUser, _ := dto.UpdateUserFromDTO(user)
	assert.Equal(t, user.Email, newUser.Email)
	assert.Equal(t, user.Password, newUser.Password)
	assert.Equal(t, user.Username, newUser.Username)
	assert.Equal(t, companyNIP, newUser.Company.Nip)
}

func TestUpdateUserFromDTO_CompanyNIPAsPerson(t *testing.T) {
	companyNIP := "new_nip"
	dto := user.UpdateUserDTO{
		ID:         1,
		CompanyNIP: &companyNIP,
	}
	user_ := createPerson(1)
	_, err := dto.UpdateUserFromDTO(user_)
	assert.ErrorIs(t, err, user.ErrUpdateCompany)
}

func TestUpdateUserFromDTO_PersonNameAsPerson(t *testing.T) {
	name := "new_name"
	dto := user.UpdateUserDTO{
		ID:         1,
		PersonName: &name,
	}
	user := createPerson(1)
	newUser, _ := dto.UpdateUserFromDTO(user)
	assert.Equal(t, user.Email, newUser.Email)
	assert.Equal(t, user.Password, newUser.Password)
	assert.Equal(t, user.Username, newUser.Username)
	assert.Equal(t, name, newUser.Person.Name)
}

func TestUpdateUserFromDTO_PersonNameAsCompany(t *testing.T) {
	name := "new_name"
	dto := user.UpdateUserDTO{
		ID:         1,
		PersonName: &name,
	}
	user_ := createCompany(1)
	_, err := dto.UpdateUserFromDTO(user_)
	assert.ErrorIs(t, err, user.ErrUpdatePerson)
}

func TestUpdateUserFromDTO_PersonSurnameAsPerson(t *testing.T) {
	surname := "new_surname"
	dto := user.UpdateUserDTO{
		ID:            1,
		PersonSurname: &surname,
	}
	user := createPerson(1)
	newUser, _ := dto.UpdateUserFromDTO(user)
	assert.Equal(t, user.Email, newUser.Email)
	assert.Equal(t, user.Password, newUser.Password)
	assert.Equal(t, user.Username, newUser.Username)
	assert.Equal(t, surname, newUser.Person.Surname)
}

func TestUpdateUserFromDTO_PersonSurnameAsCompany(t *testing.T) {
	surname := "new_surname"
	dto := user.UpdateUserDTO{
		ID:            1,
		PersonSurname: &surname,
	}
	user_ := createCompany(1)
	_, err := dto.UpdateUserFromDTO(user_)
	assert.ErrorIs(t, err, user.ErrUpdatePerson)
}

func TestUpdateUserFromDTO_Empty(t *testing.T) {
	dto := user.UpdateUserDTO{
		ID: 1,
	}
	user := createPerson(1)
	newUser, _ := dto.UpdateUserFromDTO(user)
	assert.Equal(t, user.Email, newUser.Email)
	assert.Equal(t, user.Password, newUser.Password)
	assert.Equal(t, user.Username, newUser.Username)
	assert.Equal(t, user.Selector, newUser.Selector)
}
