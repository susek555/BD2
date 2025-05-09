package review

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"gorm.io/gorm"
)

type ReviewServiceInterface interface {
	generic.CRUDService[Review]
	GetByReviewerId(reviewerId uint) ([]Review, error)
	GetByReviewedId(reviewedId uint) ([]Review, error)
	GetByReviewerIdAndReviewedId(reviewerId uint, reviewedId uint) (Review, error)
}

type ReviewService struct {
	generic.GenericService[Review, *ReviewRepository]
}

func NewReviewService(db *gorm.DB) *ReviewService {
	repo := NewReviewRepository(db)
	return &ReviewService{
		GenericService: generic.GenericService[Review, *ReviewRepository]{
			Repo: repo,
		},
	}
}

func (service *ReviewService) GetByReviewerId(reviewerId uint) ([]Review, error) {
	return service.Repo.GetByReviewerId(reviewerId)
}

func (service *ReviewService) GetByReviewedId(reviewedId uint) ([]Review, error) {
	return service.Repo.GetByReviewedId(reviewedId)
}

func (service *ReviewService) GetByReviewerAndReviewedId(reviewerId uint, reviewedId uint) (Review, error) {
	return service.GetByReviewerAndReviewedId(reviewerId, reviewedId)
}
