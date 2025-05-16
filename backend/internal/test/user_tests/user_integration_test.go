package user_tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	u "github.com/susek555/BD2/car-dealer-api/internal/test/test_utils"
)

// -------------------
// Get all users tests
// -------------------

func TestGetAllUsers_EmptyDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []user.User{}
	server, _ := newTestServer(seedUsers)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/users/", nil, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got []user.RetrieveUserDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Empty(t, got)
}

func TestGetAllUsers_SinglePerson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []user.User{*createPerson(1)}
	server, _ := newTestServer(seedUsers)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/users/", nil, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got []user.RetrieveUserDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(got), len(seedUsers))
	assert.Equal(t, doesUserAndRetrieveUserDTOsMatch(seedUsers[0], got[0]), true)
}

func TestGetAllUsers_SingleCopany(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []user.User{*createCompany(1)}
	server, _ := newTestServer(seedUsers)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/users/", nil, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got []user.RetrieveUserDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(got), len(seedUsers))
	assert.Equal(t, doesUserAndRetrieveUserDTOsMatch(seedUsers[0], got[0]), true)
}

func TestGetAllUsers_MultiplePeople(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []user.User{*createPerson(1), *createPerson(2)}
	server, err := newTestServer(seedUsers)
	assert.NoError(t, err)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/users/", nil, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got []user.RetrieveUserDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(got), len(seedUsers))
	for i := range len(seedUsers) {
		assert.Equal(t, doesUserAndRetrieveUserDTOsMatch(seedUsers[i], got[i]), true)
	}
}

func TestGetAllUsers_MultipleCompanies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []user.User{*createCompany(1), *createCompany(2)}
	server, err := newTestServer(seedUsers)
	assert.NoError(t, err)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/users/", nil, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got []user.RetrieveUserDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(got), len(seedUsers))
	for i := range len(seedUsers) {
		assert.Equal(t, doesUserAndRetrieveUserDTOsMatch(seedUsers[i], got[i]), true)
	}
}

func TestGetAllUsers_Mixed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []user.User{*createPerson(1), *createCompany(2), *createPerson(3), *createCompany(4)}
	server, err := newTestServer(seedUsers)
	assert.NoError(t, err)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/users/", nil, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got []user.RetrieveUserDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(got), len(seedUsers))
	for i := range len(seedUsers) {
		assert.Equal(t, doesUserAndRetrieveUserDTOsMatch(seedUsers[i], got[i]), true)
	}
}
