package review_tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/review"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ------
// Setup
// ------

func setupDB(users []models.User, reviews []models.Review) (review.ReviewRepositoryInterface, user.UserRepositoryInterface, error) {
	dsn := "host=localhost user=bd2_user password=bd2_password dbname=bd2_test port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	db.Exec("TRUNCATE TABLE reviews, companies, people, users RESTART IDENTITY CASCADE")
	userRepo := user.NewUserRepository(db)
	for _, u := range users {
		err = userRepo.Create(&u)
		if err != nil {
			return nil, nil, err
		}
	}
	reviewRepo := review.NewReviewRepository(db)
	for _, r := range reviews {
		err = reviewRepo.Create(&r)
		if err != nil {
			return nil, nil, err
		}
	}
	return reviewRepo, userRepo, nil
}

func newTestServer(seedUsers []models.User, seedReviews []models.Review) (*gin.Engine, review.ReviewServiceInterface, user.UserServiceInterface, error) {
	reviewRepo, userRepo, err := setupDB(seedUsers, seedReviews)
	if err != nil {
		return nil, nil, nil, err
	}
	verifier := jwt.NewJWTVerifier("secret")
	reviewService := review.NewReviewService(reviewRepo)
	userService := user.NewUserService(userRepo)
	reviewHandler := review.NewHandler(reviewService)

	router := gin.Default()
	reviewRoutes := router.Group("/review")
	reviewRoutes.GET("/", reviewHandler.GetAllReviews)
	reviewRoutes.GET("/:id", reviewHandler.GetReviewById)
	reviewRoutes.POST("/", middleware.Authenticate(verifier), reviewHandler.CreateReview)
	reviewRoutes.PUT("/", middleware.Authenticate(verifier), reviewHandler.UpdateReview)
	reviewRoutes.DELETE("/:id", middleware.Authenticate(verifier), reviewHandler.DeleteReview)
	reviewRoutes.POST("/reviewer/:id", reviewHandler.GetReviewsByReviewerId)
	reviewRoutes.POST("/reviewee/:id", reviewHandler.GetReviewsByRevieweeId)
	reviewRoutes.GET("/reviewer/reviewee/:reviewerId/:revieweeId", reviewHandler.GetReviewsByReviewerIdAndRevieweeId)

	return router, reviewService, userService, nil
}

func getValidToken(userId uint, email string) (string, error) {
	secret := []byte("secret")
	return jwt.GenerateToken(email, int64(userId), secret, time.Now().Add(1*time.Hour))
}

func TestGetAllReviewsNoReviews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedUsers []models.User
	var seedReviews []models.Review
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusOK
	req := httptest.NewRequest(http.MethodGet, "/review/", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got []review.RetrieveReviewDTO
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Empty(t, got)
}

func TestGetAllReviewsOneReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "unique@example.com",
			Username: "taken_username",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "unique2@example.com",
			Username: "taken_username2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusOK
	req := httptest.NewRequest(http.MethodGet, "/review/", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got []review.RetrieveReviewDTO
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, seedReviews[0].Rating, got[0].Rating)
	assert.Equal(t, seedReviews[0].Description, got[0].Description)
	assert.Equal(t, seedReviews[0].ReviewerID, got[0].Reviewer.ID)
	assert.Equal(t, seedReviews[0].RevieweeId, got[0].Reviewee.ID)
	assert.Equal(t, uint(1), got[0].ID)
}

func TestGetAllReviewsMultipleReviews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles3@gmail.com",
			Username: "herakles3",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
		{
			ReviewerID:  1,
			RevieweeId:  3,
			Rating:      4,
			Description: "Good service!",
		},
		{
			ReviewerID:  2,
			RevieweeId:  1,
			Rating:      3,
			Description: "Average service!",
		},
		{
			ReviewerID:  2,
			RevieweeId:  3,
			Rating:      2,
			Description: "Bad service!",
		},
		{
			ReviewerID:  3,
			RevieweeId:  1,
			Rating:      1,
			Description: "Terrible service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusOK
	req := httptest.NewRequest(http.MethodGet, "/review/", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got []review.RetrieveReviewDTO
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Len(t, got, len(seedReviews))
	for i, review_ := range seedReviews {
		assert.Equal(t, review_.Rating, got[i].Rating)
		assert.Equal(t, review_.Description, got[i].Description)
		assert.Equal(t, review_.ReviewerID, got[i].Reviewer.ID)
		assert.Equal(t, review_.RevieweeId, got[i].Reviewee.ID)
		assert.Equal(t, uint(i+1), got[i].ID)
	}
}

func TestGetReviewByIdNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles1@gmail.com",
			Username: "herakles1",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusNotFound
	req := httptest.NewRequest(http.MethodGet, "/review/999", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got review.RetrieveReviewDTO
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Empty(t, got)
}

func TestGetReviewById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusOK
	req := httptest.NewRequest(http.MethodGet, "/review/1", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got review.RetrieveReviewDTO
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, seedReviews[0].Rating, got.Rating)
	assert.Equal(t, seedReviews[0].Description, got.Description)
	assert.Equal(t, seedReviews[0].ReviewerID, got.Reviewer.ID)
	assert.Equal(t, seedReviews[0].RevieweeId, got.Reviewee.ID)
	assert.Equal(t, uint(1), got.ID)
}

func TestCreateReviewNoAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedReviews []models.Review
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusUnauthorized
	req := httptest.NewRequest(http.MethodPost, "/review/", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "unauthorized", got["message"])
}

func TestCreateReviewInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedReviews []models.Review
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusForbidden
	req := httptest.NewRequest(http.MethodPost, "/review/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "forbidden", got["message"])
}

func TestCreateReviewSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedReviews []models.Review
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusCreated
	reviewInput := review.CreateReviewDTO{
		Rating:      5,
		Description: "Great service!",
		RevieweeId:  2,
	}
	reviewInputJSON, err := json.Marshal(reviewInput)
	assert.NoError(t, err)
	token, err := getValidToken(1, seedUsers[0].Email)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/review/", strings.NewReader(string(reviewInputJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got review.RetrieveReviewDTO
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, reviewInput.Rating, got.Rating)
	assert.Equal(t, reviewInput.Description, got.Description)
	assert.Equal(t, reviewInput.RevieweeId, got.Reviewee.ID)
	assert.Equal(t, uint(1), got.Reviewer.ID)
	assert.Equal(t, uint(1), got.ID)
}

func TestCreateReviewInvalidRating(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedReviews []models.Review
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusBadRequest
	reviewInput := review.CreateReviewDTO{
		Rating:      6,
		Description: "Great service!",
		RevieweeId:  2,
	}
	reviewInputJSON, err := json.Marshal(reviewInput)
	assert.NoError(t, err)
	token, err := getValidToken(1, seedUsers[0].Email)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/review/", strings.NewReader(string(reviewInputJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "invalid rating, must be between 1 and 5", got["error_description"])
}

func TestCreateReviewSelfReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedReviews []models.Review
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusBadRequest
	reviewInput := review.CreateReviewDTO{
		Rating:      5,
		Description: "Great service!",
		RevieweeId:  1,
	}
	reviewInputJSON, err := json.Marshal(reviewInput)
	assert.NoError(t, err)
	token, err := getValidToken(1, seedUsers[0].Email)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/review/", strings.NewReader(string(reviewInputJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "CHECK constraint failed: chk_reviews_reviewer_id", got["error_description"])
}

func TestCreateReviewReviewAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusBadRequest
	reviewInput := review.CreateReviewDTO{
		Rating:      5,
		Description: "Great service!",
		RevieweeId:  2,
	}
	reviewInputJSON, err := json.Marshal(reviewInput)
	assert.NoError(t, err)
	token, err := getValidToken(1, seedUsers[0].Email)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/review/", strings.NewReader(string(reviewInputJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "UNIQUE constraint failed: reviews.reviewer_id, reviews.reviewee_id", got["error_description"])
}

func TestUpdateReviewNoAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusUnauthorized
	req := httptest.NewRequest(http.MethodPut, "/review/", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "unauthorized", got["message"])
}

func TestUpdateReviewInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusForbidden
	req := httptest.NewRequest(http.MethodPut, "/review/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "forbidden", got["message"])
}

func TestUpdateReviewSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusOK
	reviewInput := review.UpdateReviewDTO{
		ID:          1,
		Rating:      4,
		Description: "Good service!",
	}
	reviewInputJSON, err := json.Marshal(reviewInput)
	assert.NoError(t, err)
	token, err := getValidToken(1, seedUsers[0].Email)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPut, "/review/", strings.NewReader(string(reviewInputJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got review.RetrieveReviewDTO
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, reviewInput.Rating, got.Rating)
	assert.Equal(t, reviewInput.Description, got.Description)
	assert.Equal(t, seedReviews[0].ReviewerID, got.Reviewer.ID)
	assert.Equal(t, seedReviews[0].RevieweeId, got.Reviewee.ID)
	assert.Equal(t, uint(1), got.ID)
	assert.Equal(t, uint(1), got.Reviewer.ID)
	assert.Equal(t, uint(2), got.Reviewee.ID)
}

func TestUpdateReviewNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusBadRequest
	reviewInput := review.UpdateReviewDTO{
		ID:          999,
		Rating:      4,
		Description: "Good service!",
	}
	reviewInputJSON, err := json.Marshal(reviewInput)
	assert.NoError(t, err)
	token, err := getValidToken(1, seedUsers[0].Email)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPut, "/review/", strings.NewReader(string(reviewInputJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "record not found", got["error_description"])
}

func TestUpdateReviewInvalidRating(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusBadRequest
	reviewInput := review.UpdateReviewDTO{
		ID:          1,
		Rating:      6,
		Description: "Good service!",
	}
	reviewInputJSON, err := json.Marshal(reviewInput)
	assert.NoError(t, err)
	token, err := getValidToken(1, seedUsers[0].Email)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPut, "/review/", strings.NewReader(string(reviewInputJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "invalid rating, must be between 1 and 5", got["error_description"])
}

func TestUpdateNotYourReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles3@gmail.com",
			Username: "herakles3",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusBadRequest
	// User 3 tries to update a review made by user 1
	// to user 2
	// User 3 should not be able to do that
	// because he is not the reviewer
	reviewInput := review.UpdateReviewDTO{
		ID:          1,
		Rating:      4,
		Description: "Good service!",
	}
	reviewInputJSON, err := json.Marshal(reviewInput)
	assert.NoError(t, err)
	token, err := getValidToken(3, seedUsers[2].Email)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPut, "/review/", strings.NewReader(string(reviewInputJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "you are not the reviewer of this review", got["error_description"])
}

func TestDeleteReviewNoAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusUnauthorized
	req := httptest.NewRequest(http.MethodDelete, "/review/1", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "unauthorized", got["message"])
}

func TestDeleteReviewInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusForbidden
	req := httptest.NewRequest(http.MethodDelete, "/review/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "forbidden", got["message"])
}

func TestDeleteReviewSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusNoContent
	req := httptest.NewRequest(http.MethodDelete, "/review/1", nil)
	req.Header.Set("Content-Type", "application/json")
	token, err := getValidToken(1, seedUsers[0].Email)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	// Check if the review was actually deleted
	req = httptest.NewRequest(http.MethodGet, "/review/1", nil)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "no review found", got["error_description"])
}

func TestDeleteReviewNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusBadRequest
	req := httptest.NewRequest(http.MethodDelete, "/review/999", nil)
	req.Header.Set("Content-Type", "application/json")
	token, err := getValidToken(1, seedUsers[0].Email)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "record not found", got["error_description"])
}

func TestDeleteNotYourReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusBadRequest
	// User 2 tries to delete a review made by user 1
	// to user 2
	// User 2 should not be able to do that
	// because he is not the reviewer
	req := httptest.NewRequest(http.MethodDelete, "/review/1", nil)
	req.Header.Set("Content-Type", "application/json")
	token, err := getValidToken(2, seedUsers[1].Email)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "you are not the reviewer of this review", got["error_description"])
}

func TestGetReviewsByReviewerIdNoReviews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herakles",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herakles2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedReviews []models.Review
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusOK
	input := `
	{
    	"pagination": {"page": 1, "page_size": 4}
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/review/reviewer/1", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got review.RetrieveReviewsWithPagination
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Len(t, got.Reviews, 0)
	assert.Equal(t, int64(0), got.PaginationResponse.TotalRecords)
	assert.Equal(t, int64(1), got.PaginationResponse.TotalPages)
}

func TestGetReviewsByReviewerId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herkules2@gmail.com",
			Username: "herkules2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
		{
			ReviewerID:  2,
			RevieweeId:  1,
			Rating:      4,
			Description: "Good service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusOK
	input := `
	{
		"pagination": {"page": 1, "page_size": 4}
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/review/reviewer/1", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got review.RetrieveReviewsWithPagination
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Len(t, got.Reviews, 1)
	assert.Equal(t, seedReviews[0].Rating, got.Reviews[0].Rating)
	assert.Equal(t, seedReviews[0].Description, got.Reviews[0].Description)
	assert.Equal(t, seedReviews[0].ReviewerID, got.Reviews[0].Reviewer.ID)
	assert.Equal(t, seedReviews[0].RevieweeId, got.Reviews[0].Reviewee.ID)
	assert.Equal(t, int64(1), got.PaginationResponse.TotalRecords)
	assert.Equal(t, int64(1), got.PaginationResponse.TotalPages)
}

func TestGetReviewsByReviewerIdNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herkules2@gmail.com",
			Username: "herkules2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
		{
			ReviewerID:  2,
			RevieweeId:  1,
			Rating:      4,
			Description: "Good service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusOK
	input := `
	{
		"pagination": {"page": 1, "page_size": 4}
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/review/reviewer/999", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got review.RetrieveReviewsWithPagination
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Len(t, got.Reviews, 0)
	assert.Equal(t, int64(0), got.PaginationResponse.TotalRecords)
	assert.Equal(t, int64(1), got.PaginationResponse.TotalPages)
}

func TestGetReviewsByRevieweeIdNoReviews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herkules2@gmail.com",
			Username: "herkules2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
	}
	var seedReviews []models.Review
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusOK
	input := `
	{
		"pagination": {"page": 1, "page_size": 4}
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/review/reviewee/1", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got review.RetrieveReviewsWithPagination
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Len(t, got.Reviews, 0)
	assert.Equal(t, int64(0), got.PaginationResponse.TotalRecords)
	assert.Equal(t, int64(1), got.PaginationResponse.TotalPages)
}

func TestGetReviewsByRevieweeId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herkules2@gmail.com",
			Username: "herkules2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
		{
			ReviewerID:  2,
			RevieweeId:  1,
			Rating:      4,
			Description: "Good service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusOK
	input := `
	{
		"pagination": {"page": 1, "page_size": 4}
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/review/reviewee/2", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got review.RetrieveReviewsWithPagination
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Len(t, got.Reviews, 1)
	assert.Equal(t, seedReviews[0].Rating, got.Reviews[0].Rating)
	assert.Equal(t, seedReviews[0].Description, got.Reviews[0].Description)
	assert.Equal(t, seedReviews[0].ReviewerID, got.Reviews[0].Reviewer.ID)
	assert.Equal(t, seedReviews[0].RevieweeId, got.Reviews[0].Reviewee.ID)
	assert.Equal(t, int64(1), got.PaginationResponse.TotalRecords)
	assert.Equal(t, int64(1), got.PaginationResponse.TotalPages)
}

func TestGetReviewsByRevieweeIdNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herkules2@gmail.com",
			Username: "herkules2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
		{
			ReviewerID:  2,
			RevieweeId:  1,
			Rating:      4,
			Description: "Good service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusOK
	input := `
	{
		"pagination": {"page": 1, "page_size": 4}
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/review/reviewee/999", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got review.RetrieveReviewsWithPagination
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Len(t, got.Reviews, 0)
	assert.Equal(t, int64(0), got.PaginationResponse.TotalRecords)
	assert.Equal(t, int64(1), got.PaginationResponse.TotalPages)
}

func TestGetReviewsByReviewerIdAndRevieweeIdNoReviews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herkules2@gmail.com",
			Username: "herkules2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
	}
	var seedReviews []models.Review
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusNotFound
	req := httptest.NewRequest(http.MethodGet, "/review/reviewer/1/reviewee/2", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
}

func TestGetReviewsByReviewerIdAndRevieweeId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herkules2@gmail.com",
			Username: "herkules2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herkules",
				Surname: "Wielki",
			},
		},
	}
	seedReviews := []models.Review{
		{
			ReviewerID:  1,
			RevieweeId:  2,
			Rating:      5,
			Description: "Great service!",
		},
		{
			ReviewerID:  2,
			RevieweeId:  1,
			Rating:      4,
			Description: "Good service!",
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedReviews)
	assert.NoError(t, err)
	wantStatus := http.StatusOK
	req := httptest.NewRequest(http.MethodGet, "/review/reviewer/reviewee/1/2", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got review.RetrieveReviewDTO
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, seedReviews[0].Rating, got.Rating)
	assert.Equal(t, seedReviews[0].Description, got.Description)
	assert.Equal(t, seedReviews[0].ReviewerID, got.Reviewer.ID)
	assert.Equal(t, seedReviews[0].RevieweeId, got.Reviewee.ID)
	assert.Equal(t, uint(1), got.ID)
	assert.Equal(t, uint(1), got.Reviewer.ID)
	assert.Equal(t, uint(2), got.Reviewee.ID)
	assert.Equal(t, seedUsers[0].Username, got.Reviewer.Username)
	assert.Equal(t, seedUsers[1].Username, got.Reviewee.Username)
}
