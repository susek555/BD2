package review

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"gorm.io/gorm"
)

type ReviewRepositoryInterface interface {
	generic.CRUDRepository[Review]
	GetByReviewerId(reviewerId uint) ([]Review, error)
	GetByReviewedId(reviewedId uint) ([]Review, error)
	GetByReviewerIdAndReviewedId(reviewerId uint, reviewedId uint) (Review, error)
}

type ReviewRepository struct {
	repository *generic.GormRepository[Review]
}

func NewReviewRepository(dbHandle *gorm.DB) *ReviewRepository {
	return &ReviewRepository{repository: generic.GetGormRepository[Review](dbHandle)}
}

func (repo *ReviewRepository) Create(review *Review) error {
	return repo.repository.Create(review)
}

func (repo *ReviewRepository) GetAll() ([]Review, error) {
	return repo.repository.GetAll()
}

func (repo *ReviewRepository) GetById(id uint) (Review, error) {
	return repo.repository.GetById(id)
}

func (repo *ReviewRepository) Update(review *Review) error {
	return repo.repository.Update(review)
}

func (repo *ReviewRepository) Delete(id uint) error {
	return repo.repository.Delete(id)
}

func (repo *ReviewRepository) GetByReviewerId(reviewerId uint) ([]Review, error) {
	var reviews []Review
	err := repo.repository.
		DB.
		Where("reviewer_id = ?", reviewerId).
		Preload("Reviewed").
		Preload("Reviewer").
		Find(&reviews).
		Error
	return reviews, err
}

func (repo *ReviewRepository) GetByReviewedId(reviewedId uint) ([]Review, error) {
	var reviews []Review
	err := repo.repository.
		DB.
		Where("reviewer_id = ?", reviewedId).
		Preload("Reviewed").
		Preload("Reviewer").
		Find(&reviews).
		Error
	return reviews, err
}

func (repo *ReviewRepository) GetByReviewerIdAndReviewedId(reviewerId uint, reviewedId uint) (Review, error) {
	var review Review
	err := repo.repository.
		DB.
		Where("reviewer_id = ?", reviewerId).
		Where("reviewed_id = ?", reviewedId).
		Preload("Reviewer").
		Preload("Reviewed").
		First(&review).
		Error
	return review, err
}
