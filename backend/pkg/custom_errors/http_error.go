package custom_errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Description string `json:"error_description"`
}

func NewHTTPError(description string) *HTTPError {
	return &HTTPError{Description: description}
}

func GetStatusCode(err error, errorMap map[error]int) int {
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
