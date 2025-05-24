package custom_errors

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505": // unique_violation
			return http.StatusConflict
		case "23503": // foreign_key_violation
			return http.StatusBadRequest
		case "23502": // not_null_violation
			return http.StatusBadRequest
		case "23514": // check_violation
			return http.StatusBadRequest
		case "42703": // undefined_column
			return http.StatusBadRequest
		case "42P01": // undefined_table
			return http.StatusInternalServerError
		}
	}
	if statusCode, ok := errorMap[err]; ok {
		return statusCode
	}
	return http.StatusBadRequest
}

func HandleError(c *gin.Context, err error, errorMap map[error]int) {
	code := GetStatusCode(err, errorMap)
	errorMessage := normalizeErrorMessage(err)
	httpError := HTTPError{Description: errorMessage}
	c.JSON(code, httpError)
}

func normalizeErrorMessage(err error) string {
	errorMsg := err.Error()

	// Handle postgres constraint errors to match sqlite format
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			if strings.Contains(errorMsg, "ux_review_pair") {
				return "UNIQUE constraint failed: reviews.reviewer_id, reviews.reviewee_id"
			}
			return errorMsg
		case "23514":
			if strings.Contains(errorMsg, "chk_reviews_reviewer_id") {
				return "CHECK constraint failed: chk_reviews_reviewer_id"
			}
			return errorMsg
		}
	}

	return errorMsg
}
