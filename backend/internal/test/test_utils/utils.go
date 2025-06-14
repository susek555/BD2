package test_utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
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

const JWTSECRET = "secret"

func GetValidToken(userID uint, email string) (string, error) {
	secret := []byte("secret")
	return jwt.GenerateToken(email, int64(userID), secret, time.Now().Add(1*time.Hour))
}

func InsertRecordsIntoDB[T any](db *gorm.DB, records []T) error {
	repo := generic.GetGormRepository[T](db)
	for _, record := range records {
		if err := repo.Create(&record); err != nil {
			return err
		}
	}
	return nil
}

func CloseDBConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get database connection")
	}
	if err := sqlDB.Close(); err != nil {
		panic("Failed to close database connection")
	}
}

func CleanDB(db *gorm.DB) error {
	if err := db.Exec("TRUNCATE TABLE bids, liked_offers, sale_offers, auctions RESTART IDENTITY CASCADE").Error; err != nil {
		return err
	}
	return nil
}

func GetTestCloudinary() *cloudinary.Cloudinary {
	os.Setenv("CLOUDINARY_URL", "cloudinary://232679685254738:xU13-pBToKDG825l_LhM47k8k8o@du9datfva")
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	cld, _ := cloudinary.NewFromURL(cloudinaryURL)
	return cld
}
