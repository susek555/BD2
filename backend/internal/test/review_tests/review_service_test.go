//go:build unit
// +build unit

package review_tests

import (
	"errors"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/review"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
)

// helper
func newServiceWithMock() (*review.ReviewService, *mocks.ReviewRepositoryInterface) {
	repoMock := new(mocks.ReviewRepositoryInterface)
	svc := &review.ReviewService{
		GenericService: generic.GenericService[review.Review, review.ReviewRepositoryInterface]{
			Repo: repoMock,
		},
	}
	return svc, repoMock
}

// -------------------------
// Custom query methods
// -------------------------

func TestGetByReviewerId_Success(t *testing.T) {
	svc, repo := newServiceWithMock()

	want := []review.Review{{ID: 1, ReviewerID: 10}, {ID: 2, ReviewerID: 10}}
	repo.On("GetByReviewerId", uint(10)).Return(want, nil).Once()

	got, err := svc.GetByReviewerId(10)

	require.NoError(t, err)
	assert.Equal(t, want, got)
	repo.AssertExpectations(t)
}

func TestGetByReviewerId_Error(t *testing.T) {
	svc, repo := newServiceWithMock()

	repoErr := errors.New("db error")
	repo.On("GetByReviewerId", uint(10)).Return(nil, repoErr).Once()

	got, err := svc.GetByReviewerId(10)

	require.ErrorIs(t, err, repoErr)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

func TestGetByRevieweeId_Success(t *testing.T) {
	svc, repo := newServiceWithMock()

	want := []review.Review{{ID: 3, RevieweeId: 20}}
	repo.On("GetByRevieweeId", uint(20)).Return(want, nil).Once()

	got, err := svc.GetByRevieweeId(20)

	require.NoError(t, err)
	assert.Equal(t, want, got)
	repo.AssertExpectations(t)
}

func TestGetByReviewerAndReviewee_Success(t *testing.T) {
	svc, repo := newServiceWithMock()

	want := review.Review{ID: 4, ReviewerID: 10, RevieweeId: 20}
	repo.On("GetByReviewerIdAndRevieweeId", uint(10), uint(20)).Return(want, nil).Once()

	got, err := svc.GetByReviewerIdAndRevieweeId(10, 20)

	require.NoError(t, err)
	assert.Equal(t, want, got)
	repo.AssertExpectations(t)
}

func TestGetByReviewerAndReviewee_Error(t *testing.T) {
	svc, repo := newServiceWithMock()

	repoErr := errors.New("not found")
	repo.On("GetByReviewerIdAndRevieweeId", uint(10), uint(20)).Return(review.Review{}, repoErr).Once()

	got, err := svc.GetByReviewerIdAndRevieweeId(10, 20)

	require.ErrorIs(t, err, repoErr)
	assert.Equal(t, review.Review{}, got)
	repo.AssertExpectations(t)
}

// CRUD

func TestCreate_Success(t *testing.T) {
	svc, repo := newServiceWithMock()

	in := &review.Review{Description: "ok", ReviewerID: 1, RevieweeId: 2}
	repo.On("Create", in).Return(nil).Once()

	err := svc.Create(in)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCreate_Error(t *testing.T) {
	svc, repo := newServiceWithMock()

	in := &review.Review{Description: "bad"}
	repoErr := errors.New("insert failed")
	repo.On("Create", in).Return(repoErr).Once()

	err := svc.Create(in)

	require.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}

func TestGet_Success(t *testing.T) {
	svc, repo := newServiceWithMock()

	want := review.Review{ID: 7, ReviewerID: 1, RevieweeId: 2}
	repo.On("GetById", uint(7)).Return(&want, nil).Once()

	got, err := svc.GetById(7)

	require.NoError(t, err)
	assert.Equal(t, &want, got)
	repo.AssertExpectations(t)
}

func TestGet_Error(t *testing.T) {
	svc, repo := newServiceWithMock()

	repoErr := errors.New("not found")
	repo.On("GetById", uint(7)).Return(&review.Review{}, repoErr).Once()

	got, err := svc.GetById(7)
	var want *review.Review = nil

	require.ErrorIs(t, err, repoErr)
	assert.Equal(t, want, got)
	repo.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	svc, repo := newServiceWithMock()

	upd := &review.Review{ID: 5, Description: "new desc"}
	repo.On("Update", upd).Return(nil).Once()

	err := svc.Update(upd)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUpdate_Error(t *testing.T) {
	svc, repo := newServiceWithMock()

	upd := &review.Review{ID: 5}
	repoErr := errors.New("update failed")
	repo.On("Update", upd).Return(repoErr).Once()

	err := svc.Update(upd)

	require.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	svc, repo := newServiceWithMock()

	repo.On("Delete", uint(9)).Return(nil).Once()

	err := svc.Delete(9)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDelete_Error(t *testing.T) {
	svc, repo := newServiceWithMock()

	repoErr := errors.New("delete failed")
	repo.On("Delete", uint(9)).Return(repoErr).Once()

	err := svc.Delete(9)

	require.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}
