package user_tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	u "github.com/susek555/BD2/car-dealer-api/internal/test/test_utils"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
	"gorm.io/gorm"
)

// -------------------
// Get all users tests
// -------------------

func TestGetAllUsers_EmptyDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedUsers []models.User
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got []user.RetrieveUserDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Empty(t, got)
}

func TestGetAllUsers_SinglePerson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got []user.RetrieveUserDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(got), len(seedUsers))
	assert.Equal(t, doUserAndRetrieveUserDTOsMatch(seedUsers[0], got[0]), true)
}

func TestGetAllUsers_SingleCompany(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got []user.RetrieveUserDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(got), len(seedUsers))
	assert.Equal(t, doUserAndRetrieveUserDTOsMatch(seedUsers[0], got[0]), true)
}

func TestGetAllUsers_MultiplePeople(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1), *createPerson(2)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got []user.RetrieveUserDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(got), len(seedUsers))
	for i := range len(seedUsers) {
		assert.Equal(t, doUserAndRetrieveUserDTOsMatch(seedUsers[i], got[i]), true)
	}
}

func TestGetAllUsers_MultipleCompanies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1), *createCompany(2)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got []user.RetrieveUserDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(got), len(seedUsers))
	for i := range len(seedUsers) {
		assert.Equal(t, doUserAndRetrieveUserDTOsMatch(seedUsers[i], got[i]), true)
	}
}

func TestGetAllUsers_Mixed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1), *createCompany(2), *createPerson(3), *createCompany(4)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got []user.RetrieveUserDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(got), len(seedUsers))
	for i := range len(seedUsers) {
		assert.Equal(t, doUserAndRetrieveUserDTOsMatch(seedUsers[i], got[i]), true)
	}
}

// --------------------
// Get by user id tests
// --------------------

func TestGetUserByID_EmptyDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedUsers []models.User
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/id/1", nil, nil)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
}

func TestGetByUserID_NonExistentID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/id/2", nil, nil)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
}

func TestGetUserByID_NegativeID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/id/-1", nil, nil)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
}

func TestGetUserByID_StringID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/id/abc", nil, nil)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
}

func TestGetUserByID_Person(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/id/1", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got user.RetrieveUserDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, doUserAndRetrieveUserDTOsMatch(seedUsers[0], got), true)
}

func TestGetUserByID_Company(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/id/1", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got user.RetrieveUserDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, doUserAndRetrieveUserDTOsMatch(seedUsers[0], got), true)
}

// -----------------------
// Get by user email tests
// -----------------------

func TestGetUserByEmail_EmptyDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedUsers []models.User
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/email/john1@gmail.com", nil, nil)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
}

func TestGetUserByEmail_NonExistentEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/email/john99@gmail.com", nil, nil)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
}

func TestGetUserByEmail_Person(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/email/john1@gmail.com", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got user.RetrieveUserDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, doUserAndRetrieveUserDTOsMatch(seedUsers[0], got), true)
}

func TestGetUserByEmail_Company(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/users/email/john1@gmail.com", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got user.RetrieveUserDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, doUserAndRetrieveUserDTOsMatch(seedUsers[0], got), true)
}

// -----------------
// Update user tests
// -----------------

func TestUpdateUser_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	body, err := json.Marshal(user.UpdateUserDTO{ID: 1})
	assert.NoError(t, err)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
}

func TestUpdateUser_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1), *createPerson(2)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	body, err := json.Marshal(user.UpdateUserDTO{ID: 2})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ErrInvalidUserID.Error())
}

func TestUpdateUser_UpdateUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	newUsername := "new_username"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, Username: &newUsername})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	userRepo, _, _, _, _ := getRepositories(db)
	got, err := userRepo.GetByID(seedUsers[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, got.Username, newUsername)
}

func TestUpdateUser_NotUniqueUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)
	newUsername := "new_username"
	seedUsers := []models.User{
		*createPerson(1),
		*u.Build(createPerson(2), u.WithField[models.User]("Username", newUsername)),
	}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, Username: &newUsername})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got user.UpdateResponse
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, &got.Errors)
	assert.NotEmpty(t, got.Errors["username"])
	assert.Equal(t, got.Errors["username"][0], user.ErrUsernameTaken.Error())
}

func TestUpdateUser_UpdateEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	newEmail := "newEmail@gmail.com"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, Email: &newEmail})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	userRepo, _, _, _, _ := getRepositories(db)
	got, err := userRepo.GetByID(seedUsers[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, got.Email, newEmail)
}

func TestUpdateUser_NotUniqueEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	newEmail := "newEmail@gmail.com"
	seedUsers := []models.User{
		*createPerson(1),
		*u.Build(createPerson(2), u.WithField[models.User]("Email", newEmail)),
	}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, Email: &newEmail})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got user.UpdateResponse
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, &got.Errors)
	assert.NotEmpty(t, got.Errors["email"])
	assert.Equal(t, got.Errors["email"][0], user.ErrEmailTaken.Error())
}

func TestUpdateUser_UpdatePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	newPassword := "new_password"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, Password: &newPassword})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	userRepo, _, _, _, _ := getRepositories(db)
	got, err := userRepo.GetByID(seedUsers[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, passwords.Match(newPassword, got.Password), true)
}

func TestUpdateUser_UpdateCompanyNameAsCompany(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	newCompanyName := "new_company_name"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, CompanyName: &newCompanyName})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	_, companyRepo, _, _, _ := getRepositories(db)
	got, err := companyRepo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, newCompanyName, got.Name)
}

func TestUpdateUser_UpdateCompanyNameAsPerson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	newCompanyName := "new_company_name"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, CompanyName: &newCompanyName})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got user.UpdateResponse
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Errors["other"][0], user.ErrUpdateCompany.Error())
}

func TestUpdateUser_UpdateNIPAsCompany(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	newNIP := "1234567890"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, CompanyNIP: &newNIP})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	_, companyRepo, _, _, _ := getRepositories(db)
	got, err := companyRepo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, got.Nip, newNIP)
}

func TestUpdateUser_UpdateNIPAsCompanyNotUnique(t *testing.T) {
	gin.SetMode(gin.TestMode)
	newNIP := "1234567890"
	seedUsers := []models.User{
		*createCompany(1),
		*u.Build(createCompany(2), withCompanyField(u.WithField[models.Company]("Nip", newNIP))),
	}
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, CompanyNIP: &newNIP})
	assert.NoError(t, err)
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got user.UpdateResponse
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, &got.Errors)
	assert.NotEmpty(t, got.Errors["company_nip"])
	assert.Equal(t, got.Errors["company_nip"][0], user.ErrNipAlreadyTaken.Error())
}

func TestUpdateUser_UpdateNIPAsPerson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	newNIP := "1234567890"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, CompanyNIP: &newNIP})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got user.UpdateResponse
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Errors["other"][0], user.ErrUpdateCompany.Error())
}

func TestUpdateUser_UpdatePersonNameAsPerson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	newName := "new_name"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, PersonName: &newName})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	userRepo, _, _, _, _ := getRepositories(db)
	got, err := userRepo.GetByID(seedUsers[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, got.Person.Name, newName)
}

func TestUpdateUser_UpdatePersonNameAsCompany(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	newName := "new_name"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, PersonName: &newName})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got user.UpdateResponse
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, &got.Errors)
	assert.NotEmpty(t, got.Errors["other"])
	assert.Equal(t, got.Errors["other"][0], user.ErrUpdatePerson.Error())
}

func TestUpdateUser_UpdatePersonSurnameAsPerson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	newSurname := "new_surname"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, PersonSurname: &newSurname})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	userRepo, _, _, _, _ := getRepositories(db)
	got, err := userRepo.GetByID(seedUsers[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, got.Person.Surname, newSurname)
}

func TestUpdateUser_UpdatePersonSurnameAsCompany(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	newSurname := "new_surname"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, PersonSurname: &newSurname})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got user.UpdateResponse
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, &got.Errors)
	assert.NotEmpty(t, got.Errors["other"])
	assert.Equal(t, got.Errors["other"][0], user.ErrUpdatePerson.Error())
}

// -----------------
// Delete user tests
// -----------------
func TestDeleteUser_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/users/id/1", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
}

func TestDeleteUser_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1), *createPerson(2)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/users/id/2", nil, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ErrInvalidUserID.Error())
}

func TestDeleteUser_Person(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/users/id/1", nil, &token)
	assert.Equal(t, http.StatusNoContent, receivedStatus)
	userRepo, _, personRepo, _, _ := getRepositories(db)
	_, err := userRepo.GetByID(seedUsers[0].ID)
	assert.Error(t, err, gorm.ErrRecordNotFound)
	_, err = personRepo.GetByID(seedUsers[0].ID)
	assert.Error(t, err, gorm.ErrRecordNotFound)
}

func TestDeleteUser_Company(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, nil, nil)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/users/id/1", nil, &token)
	assert.Equal(t, http.StatusNoContent, receivedStatus)
	userRepo, _, companyRepo, _, _ := getRepositories(db)
	_, err := userRepo.GetByID(seedUsers[0].ID)
	assert.Error(t, err, gorm.ErrRecordNotFound)
	_, err = companyRepo.GetByID(seedUsers[0].ID)
	assert.Error(t, err, gorm.ErrRecordNotFound)
}

func TestDeleteUser_HasUnsoldOffers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(2)}
	seedOffers := []models.SaleOffer{*createSaleOffer(1, 2)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers, seedOffers, nil)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	wantStatus := http.StatusNoContent
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/users/id/2", nil, &token)
	assert.Equal(t, wantStatus, receivedStatus)
	_, _, _, saleOfferRepo, _ := getRepositories(db)
	_, err := saleOfferRepo.GetByID(1)
	assert.Error(t, err, gorm.ErrRecordNotFound)
}

func TestDeleteUser_HasSoldOffers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(2), *createPerson(3)}
	seedOffers := []models.SaleOffer{*createSoldSaleOffer(1, 2)}
	seedPurchases := []models.Purchase{*createPurchase(1, 3)}
	db, _ := setupDBWithDeletedUser()
	server, _ := newTestServer(db, seedUsers, seedOffers, seedPurchases)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	wantStatus := http.StatusNoContent
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/users/id/2", nil, &token)
	assert.Equal(t, wantStatus, receivedStatus)
	_, _, _, saleOfferRepo, _ := getRepositories(db)
	got, err := saleOfferRepo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, got.Status, enums.SOLD)
	assert.Equal(t, got.UserID, uint(1)) // User ID of the deleted user should be set to 1 (deleted user)
}

func TestDeleteUser_HasPurchases(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(2), *createPerson(3)}
	seedOffers := []models.SaleOffer{*createSoldSaleOffer(1, 3)}
	seedPurchases := []models.Purchase{*createPurchase(1, 2)}
	db, _ := setupDBWithDeletedUser()
	server, _ := newTestServer(db, seedUsers, seedOffers, seedPurchases)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	wantStatus := http.StatusNoContent
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/users/id/2", nil, &token)
	assert.Equal(t, wantStatus, receivedStatus)
	_, _, _, _, purchaseRepo := getRepositories(db)
	got, err := purchaseRepo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, got.BuyerID, uint(1)) // Buyer ID of the deleted user should be set to 1 (deleted user)
}
