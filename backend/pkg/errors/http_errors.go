package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func GetStatusCode(err error, errorMap map[error]int) int {
	if statusCode, ok := errorMap[err]; ok {
		return statusCode
	}
	return http.StatusInternalServerError
}

func HandleError(c *gin.Context, err error, errorMap map[error]int) {
	code := GetStatusCode(err, errorMap)
	httpError := HTTPError{Code: code, Message: err.Error()}
	c.JSON(code, httpError)
}
