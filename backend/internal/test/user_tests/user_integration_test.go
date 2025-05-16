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
