package custom_errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
)

type HTTPError struct {
	Description string `json:"error_description"`
}

func NewHTTPError(description string) *HTTPError {
	return &HTTPError{Description: description}
}

func GetStatusCode(err error, errorMap map[error]int) int {
	if err == nil {
		return http.StatusOK
	}
	// TODO Replace with POSTGRESQL error handling
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrConstraint {
			return http.StatusBadRequest
		}
	}
	if statusCode, ok := errorMap[err]; ok {
		return statusCode
	}
	return http.StatusInternalServerError
}

func HandleError(c *gin.Context, err error, errorMap map[error]int) {
	code := GetStatusCode(err, errorMap)
	httpError := HTTPError{Description: err.Error()}
	c.JSON(code, httpError)
}
