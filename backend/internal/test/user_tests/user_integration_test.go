package user_tests

import (
	"encoding/json"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
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
	seedUsers := []models.User{}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
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
	seedUsers := []models.User{}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
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
	seedUsers := []models.User{}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
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
	server, _ := newTestServer(db, seedUsers)
	body, err := json.Marshal(user.UpdateUserDTO{ID: 1})
	assert.NoError(t, err)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
}

func TestUpdateUser_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1), *createPerson(2)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	body, err := json.Marshal(user.UpdateUserDTO{ID: 2})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ErrForbidden.Error())
}

func TestUpdateUser_UpdateUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	newUsername := "new_username"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, Username: &newUsername})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	userRepo, _, _ := getRepositories(db)
	got, err := userRepo.GetById(seedUsers[0].ID)
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
	server, _ := newTestServer(db, seedUsers)
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, Username: &newUsername})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, &got.Description)
}

func TestUpdateUser_UpdateEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	newEmail := "newEmail@gmail.com"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, Email: &newEmail})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	userRepo, _, _ := getRepositories(db)
	got, err := userRepo.GetById(seedUsers[0].ID)
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
	server, _ := newTestServer(db, seedUsers)
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, Email: &newEmail})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, &got.Description)
}

func TestUpdateUser_UpdatePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	newPassword := "new_password"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, Password: &newPassword})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	userRepo, _, _ := getRepositories(db)
	got, err := userRepo.GetById(seedUsers[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, passwords.Match(newPassword, got.Password), true)
}

func TestUpdateUser_UpdateCompanyNameAsCompany(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	newCompanyName := "new_company_name"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, CompanyName: &newCompanyName})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	_, companyRepo, _ := getRepositories(db)
	got, err := companyRepo.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, got.Name, newCompanyName)
}

func TestUpdateUser_UpdateCompanyNameAsPerson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	newCompanyName := "new_company_name"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, CompanyName: &newCompanyName})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, user.ErrUpdateCompany.Error())
}

func TestUpdateUser_UpdateNIPAsCompany(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	newNIP := "1234567890"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, CompanyNIP: &newNIP})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	_, companyRepo, _ := getRepositories(db)
	got, err := companyRepo.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, got.NIP, newNIP)
}

func TestUpdateUser_UpdateNIPAsCompanyNotUnique(t *testing.T) {
	gin.SetMode(gin.TestMode)
	newNIP := "1234567890"
	seedUsers := []models.User{
		*createCompany(1),
		*u.Build(createCompany(2), withCompanyField(u.WithField[models.Company]("NIP", newNIP))),
	}
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, CompanyNIP: &newNIP})
	assert.NoError(t, err)
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, &got.Description)
}

func TestUpdateUser_UpdateNIPAsPerson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	newNIP := "1234567890"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, CompanyNIP: &newNIP})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, &got.Description)
}

func TestUpdateUser_UpdatePersonNameAsPerson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	newName := "new_name"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, PersonName: &newName})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	userRepo, _, _ := getRepositories(db)
	got, err := userRepo.GetById(seedUsers[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, got.Person.Name, newName)
}

func TestUpdateUser_UpdatePersonNameAsCompany(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	newName := "new_name"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, PersonName: &newName})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, &got.Description)
}

func TestUpdateUser_UpdatePersonSurnameAsPerson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	newSurname := "new_surname"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, PersonSurname: &newSurname})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	userRepo, _, _ := getRepositories(db)
	got, err := userRepo.GetById(seedUsers[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, got.Person.Surname, newSurname)
}

func TestUpdateUser_UpdatePersonSurnameAsCompany(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	newSurname := "new_surname"
	body, err := json.Marshal(user.UpdateUserDTO{ID: seedUsers[0].ID, PersonSurname: &newSurname})
	assert.NoError(t, err)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/users/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, &got.Description)
}

// -----------------
// Delete user tests
// -----------------
func TestDeleteUser_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/users/id/1", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
}

func TestDeleteUser_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1), *createPerson(2)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/users/id/2", nil, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ErrForbidden.Error())
}

func TestDeleteUser_Person(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createPerson(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/users/id/1", nil, &token)
	assert.Equal(t, http.StatusNoContent, receivedStatus)
	userRepo, _, personRepo := getRepositories(db)
	_, err := userRepo.GetById(seedUsers[0].ID)
	assert.Error(t, err, gorm.ErrRecordNotFound)
	_, err = personRepo.GetById(seedUsers[0].ID)
	assert.Error(t, err, gorm.ErrRecordNotFound)
}

func TestDeleteUser_Company(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{*createCompany(1)}
	db, _ := setupDB()
	server, _ := newTestServer(db, seedUsers)
	token, _ := u.GetValidToken(seedUsers[0].ID, seedUsers[0].Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/users/id/1", nil, &token)
	assert.Equal(t, http.StatusNoContent, receivedStatus)
	userRepo, _, companyRepo := getRepositories(db)
	_, err := userRepo.GetById(seedUsers[0].ID)
	assert.Error(t, err, gorm.ErrRecordNotFound)
	_, err = companyRepo.GetById(seedUsers[0].ID)
	assert.Error(t, err, gorm.ErrRecordNotFound)
}
