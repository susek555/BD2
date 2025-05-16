package review

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

//go:generate mockery --name=ReviewRepositoryInterface --output=../../test/mocks --case=snake --with-expecter
type ReviewRepositoryInterface interface {
	generic.CRUDRepository[Review]
	GetByReviewerId(reviewerId uint) ([]Review, error)
	GetByRevieweeId(reviewedId uint) ([]Review, error)
	GetByReviewerIdAndRevieweeId(reviewerId uint, reviewedId uint) (*Review, error)
	GetFiltered(filter *ReviewFilter) ([]Review, *pagination.PaginationResponse, error)
	GetAverageRatingByRevieweeId(revieweeId uint) (float64, error)
}

type ReviewRepository struct {
	repository *generic.GormRepository[Review]
}

func NewReviewRepository(dbHandle *gorm.DB) ReviewRepositoryInterface {
	return &ReviewRepository{repository: generic.GetGormRepository[Review](dbHandle)}
}

func (repo *ReviewRepository) Create(review *Review) error {
	db := repo.repository.DB

	if err := db.Create(review).Error; err != nil {
		return err
	}
	return db.
		Preload("Reviewer").
		Preload("Reviewee").
		First(review, review.ID).
		Error
}

func (repo *ReviewRepository) GetAll() ([]Review, error) {
	db := repo.repository.DB
	var reviews []Review
	err := db.
		Preload("Reviewer").
		Preload("Reviewee").
		Find(&reviews).
		Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (repo *ReviewRepository) GetById(id uint) (*Review, error) {
	db := repo.repository.DB
	var review Review
	err := db.
		Preload("Reviewer").
		Preload("Reviewee").
		First(&review, id).
		Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (repo *ReviewRepository) GetFiltered(filter *ReviewFilter) ([]Review, *pagination.PaginationResponse, error) {
	query := repo.buildBaseQuery().
		Joins("JOIN users as reviewer ON reviews.reviewer_id = reviewer.id").
		Joins("JOIN users as reviewee ON reviews.reviewee_id = reviewee.id")
	query, err := filter.ApplyReviewFilters(query)
	if err != nil {
		return nil, nil, err
	}
	reviews, paginationResponse, err := pagination.PaginateResults[Review](&filter.Pagination, query)
	if err != nil {
		return nil, nil, err
	}
	return reviews, paginationResponse, nil
}

func (repo *ReviewRepository) Update(review *Review) error {
	db := repo.repository.DB
	err := db.Save(review).Error
	if err != nil {
		return err
	}
	err = db.
		Preload("Reviewer").
		Preload("Reviewee").
		First(review, review.ID).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *ReviewRepository) Delete(id uint) error {
	return repo.repository.Delete(id)
}

func (repo *ReviewRepository) GetByReviewerId(reviewerId uint) ([]Review, error) {
	var reviews []Review
	err := repo.repository.
		DB.
		Where("reviewer_id = ?", reviewerId).
		Preload("Reviewee").
		Preload("Reviewer").
		Find(&reviews).
		Error
	return reviews, err
}

func (repo *ReviewRepository) GetByRevieweeId(reviewedId uint) ([]Review, error) {
	var reviews []Review
	err := repo.repository.
		DB.
		Where("reviewee_id = ?", reviewedId).
		Preload("Reviewee").
		Preload("Reviewer").
		Find(&reviews).
		Error
	return reviews, err
}

func (repo *ReviewRepository) GetByReviewerIdAndRevieweeId(reviewerId uint, reviewedId uint) (*Review, error) {
	var review Review
	err := repo.repository.
		DB.
		Where("reviewer_id = ?", reviewerId).
		Where("reviewee_id = ?", reviewedId).
		Preload("Reviewer").
		Preload("Reviewee").
		First(&review).
		Error
	return &review, err
}

func (repo *ReviewRepository) GetAverageRatingByRevieweeId(revieweeId uint) (float64, error) {
	var average float64
	err := repo.repository.
		DB.
		Model(&Review{}).
		Select("AVG(rating)").
		Where("reviewee_id = ?", revieweeId).
		Scan(&average).
		Error
	if err != nil {
		return 0, err
	}
	return average, nil
}

func (repo *ReviewRepository) buildBaseQuery() *gorm.DB {
	db := repo.repository.DB
	query := db.
		Preload("Reviewer").
		Preload("Reviewee")
	return query
}
