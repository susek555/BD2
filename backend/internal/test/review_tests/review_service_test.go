package review_tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/review"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
)

// helper
func newServiceWithMock() (*review.ReviewService, *mocks.ReviewRepositoryInterface) {
	repoMock := new(mocks.ReviewRepositoryInterface)
	svc := &review.ReviewService{
		Repo: repoMock,
	}
	return svc, repoMock
}

// -------------------------
// Custom query methods
// -------------------------

func TestGetByReviewerID_Success(t *testing.T) {
	svc, repo := newServiceWithMock()

	mockReviews := []models.Review{
		{ID: 1, ReviewerID: 10, RevieweeID: 1, Reviewer: &models.User{ID: 10, Username: "reviewer1"}, Reviewee: &models.User{ID: 1, Username: "reviewee1"}},
		{ID: 2, ReviewerID: 10, RevieweeID: 2, Reviewer: &models.User{ID: 10, Username: "reviewer1"}, Reviewee: &models.User{ID: 2, Username: "reviewee2"}},
	}

	var expectedDTOs []review.RetrieveReviewDTO
	for _, r := range mockReviews {
		expectedDTOs = append(expectedDTOs, *review.MapToDTO(&r))
	}

	repo.On("GetByReviewerID", uint(10)).Return(mockReviews, nil).Once()

	got, err := svc.GetByReviewerID(10)

	require.NoError(t, err)
	assert.Equal(t, expectedDTOs, got)
	repo.AssertExpectations(t)
}

func TestGetByReviewerID_Error(t *testing.T) {
	svc, repo := newServiceWithMock()

	repoErr := errors.New("db error")
	repo.On("GetByReviewerID", uint(10)).Return(nil, repoErr).Once()

	got, err := svc.GetByReviewerID(10)

	require.ErrorIs(t, err, repoErr)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

func TestGetByRevieweeID_Success(t *testing.T) {
	svc, repo := newServiceWithMock()

	mockReviews := []models.Review{
		{ID: 3, ReviewerID: 5, RevieweeID: 20, Reviewer: &models.User{ID: 5, Username: "reviewer"}, Reviewee: &models.User{ID: 20, Username: "reviewee"}},
	}

	var expectedDTOs []review.RetrieveReviewDTO
	for _, r := range mockReviews {
		expectedDTOs = append(expectedDTOs, *review.MapToDTO(&r))
	}

	repo.On("GetByRevieweeID", uint(20)).Return(mockReviews, nil).Once()

	got, err := svc.GetByRevieweeID(20)

	require.NoError(t, err)
	assert.Equal(t, expectedDTOs, got)
	repo.AssertExpectations(t)
}

func TestGetByReviewerAndReviewee_Success(t *testing.T) {
	svc, repo := newServiceWithMock()

	mockReview := &models.Review{
		ID:         4,
		ReviewerID: 10,
		RevieweeID: 20,
		Reviewer:   &models.User{ID: 10, Username: "reviewer"},
		Reviewee:   &models.User{ID: 20, Username: "reviewee"},
	}
	expectedDTO := *review.MapToDTO(mockReview)

	repo.On("GetByReviewerIDAndRevieweeID", uint(10), uint(20)).Return(mockReview, nil).Once()

	got, err := svc.GetByReviewerIDAndRevieweeID(10, 20)

	require.NoError(t, err)
	assert.Equal(t, &expectedDTO, got)
	repo.AssertExpectations(t)
}

func TestGetByReviewerAndReviewee_Error(t *testing.T) {
	svc, repo := newServiceWithMock()

	repoErr := errors.New("not found")
	repo.On("GetByReviewerIDAndRevieweeID", uint(10), uint(20)).Return(nil, repoErr).Once()

	got, err := svc.GetByReviewerIDAndRevieweeID(10, 20)

	require.ErrorIs(t, err, repoErr)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

// CRUD

// --- Create ---

func TestCreate_Success(t *testing.T) {
	svc, repo := newServiceWithMock()
	userID := uint(1)

	in := &review.CreateReviewDTO{Description: "ok", Rating: 1, RevieweeID: 2}

	repo.
		On("Create", mock.AnythingOfType("*models.Review")).
		Run(func(args mock.Arguments) {
			r := args.Get(0).(*models.Review)

			r.ID = 1
			r.Reviewer = &models.User{ID: 1, Username: "author"}
			r.Reviewee = &models.User{ID: r.RevieweeID, Username: "user"}
		}).
		Return(nil).
		Once()

	got, err := svc.Create(userID, in)

	require.NoError(t, err)
	repo.AssertExpectations(t)

	assert.Equal(t, uint(1), got.ID)
	assert.Equal(t, in.Description, got.Description)
	assert.Equal(t, in.Rating, got.Rating)
	assert.Equal(t, in.RevieweeID, got.Reviewee.ID)
}

func TestCreate_Error(t *testing.T) {
	svc, repo := newServiceWithMock()
	userID := uint(1)

	in := &review.CreateReviewDTO{Description: "bad", Rating: 1}
	repoErr := errors.New("err")

	repo.
		On("Create", mock.AnythingOfType("*models.Review")).
		Return(repoErr).
		Once()

	got, err := svc.Create(userID, in)

	require.Error(t, err)
	assert.Equal(t, repoErr.Error(), err.Error())
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

func TestGet_Success(t *testing.T) {
	svc, repo := newServiceWithMock()

	want := models.Review{
		ID:         7,
		ReviewerID: 1,
		RevieweeID: 2,
		Reviewer: &models.User{
			ID:       1,
			Username: "reviewer",
		},
		Reviewee: &models.User{
			ID:       2,
			Username: "reviewee",
		},
	}
	repo.On("GetByID", uint(7)).Return(&want, nil).Once()

	got, err := svc.GetByID(7)
	expected := review.MapToDTO(&want)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
	repo.AssertExpectations(t)
}

func TestGet_Error(t *testing.T) {
	svc, repo := newServiceWithMock()

	repoErr := errors.New("no review found")
	repo.On("GetByID", uint(7)).Return(&models.Review{}, repoErr).Once()

	got, err := svc.GetByID(7)
	var want *review.RetrieveReviewDTO = nil

	require.Error(t, err)
	assert.Equal(t, repoErr.Error(), err.Error())
	assert.Equal(t, want, got)
	repo.AssertExpectations(t)
}

// --- Update ---

func TestUpdate_Success(t *testing.T) {
	svc, repo := newServiceWithMock()

	reviewerID := uint(5)
	reviewID := uint(5)

	upd := &review.UpdateReviewDTO{
		ID:          reviewID,
		Description: "new desc",
		Rating:      4, // add whatever fields your DTO allows
	}

	// ── step 1: service asks the repo for the existing review ──
	repo.
		On("GetByID", reviewID).
		Return(&models.Review{
			ID:         reviewID,
			ReviewerID: reviewerID,
			RevieweeID: 2,
			Reviewer:   &models.User{ID: reviewerID, Username: "author"},
			Reviewee:   &models.User{ID: 2, Username: "user"},
		}, nil).
		Once()

	repo.
		On("GetByReviewerIDAndRevieweeID", reviewerID, uint(2)).
		Return(&models.Review{
			ID:         reviewID,
			ReviewerID: reviewerID,
			RevieweeID: 2,
			Reviewer:   &models.User{ID: reviewerID, Username: "author"},
			Reviewee:   &models.User{ID: 2, Username: "user"},
		}, nil).
		Once()

	// ── step 2: service calls Update(&reviewObj) ──
	repo.
		On("Update", mock.AnythingOfType("*models.Review")).
		Run(func(args mock.Arguments) {
			r := args.Get(0).(*models.Review)

			// simulate DB finishing the mutation and preloading relations
			r.Reviewer = &models.User{ID: reviewerID, Username: "author"}
			r.Reviewee = &models.User{ID: r.RevieweeID, Username: "user"}
		}).
		Return(nil).
		Once()

	got, err := svc.Update(reviewerID, upd)

	require.NoError(t, err)
	repo.AssertExpectations(t)

	assert.Equal(t, reviewID, got.ID)
	assert.Equal(t, upd.Description, got.Description)
	assert.Equal(t, upd.Rating, got.Rating)
	assert.Equal(t, uint(2), got.Reviewee.ID)
}

func TestUpdate_Error(t *testing.T) {
	svc, repo := newServiceWithMock()

	reviewerID := uint(5)
	reviewID := uint(5)
	upd := &review.UpdateReviewDTO{ID: reviewID, Description: "fail", Rating: 1}

	repoErr := errors.New("update failed")

	repo.On("GetByID", reviewID).
		Return(&models.Review{ID: reviewID, ReviewerID: reviewerID, RevieweeID: 2}, nil).
		Once()

	repo.On("GetByReviewerIDAndRevieweeID", reviewerID, uint(2)).
		Return(&models.Review{ID: reviewID, ReviewerID: reviewerID, RevieweeID: 2}, nil).
		Once()

	repo.On("Update", mock.AnythingOfType("*models.Review")).Return(repoErr).Once()

	got, err := svc.Update(reviewerID, upd)

	require.ErrorIs(t, err, repoErr)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

// --- Delete ---

func TestDelete_Success(t *testing.T) {
	svc, repo := newServiceWithMock()

	userID := uint(9)
	reviewID := uint(10)

	repo.On("GetByID", reviewID).
		Return(&models.Review{ID: reviewID, ReviewerID: userID}, nil).
		Once()
	repo.On("Delete", reviewID).Return(nil).Once()

	err := svc.Delete(userID, reviewID)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDelete_Error(t *testing.T) {
	svc, repo := newServiceWithMock()

	userID := uint(9)
	reviewID := uint(10)
	repoErr := errors.New("delete failed")

	repo.On("GetByID", reviewID).
		Return(&models.Review{ID: reviewID, ReviewerID: userID}, nil).
		Once()
	repo.On("Delete", reviewID).Return(repoErr).Once()

	err := svc.Delete(userID, reviewID)

	require.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}
