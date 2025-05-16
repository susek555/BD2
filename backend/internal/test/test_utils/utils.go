package test_utils

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
)

type Option[T any] func(*T)

func WithField[T any](fieldName string, fieldValue interface{}) Option[T] {
	return func(obj *T) {
		v := reflect.ValueOf(obj).Elem()
		field := v.FieldByName(fieldName)
		field.Set(reflect.ValueOf(fieldValue))
	}
}

func Build[T any](obj *T, options ...Option[T]) *T {
	for _, option := range options {
		option(obj)
	}
	return obj
}

func PerformRequest(server *gin.Engine, method string, url string, body []byte, authToken *string) ([]byte, int) {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, url, strings.NewReader(string(body)))
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	if authToken != nil {
		req.Header.Set("Authorization", "Bearer "+*authToken)
	}
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	return w.Body.Bytes(), w.Code
}

func GetDefaultPaginationRequest() *pagination.PaginationRequest {
	return &pagination.PaginationRequest{Page: 1, PageSize: 8}
}
