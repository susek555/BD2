package review

import (
	"errors"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"math"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

//go:generate mockery --name=ReviewRepositoryInterface --output=../../test/mocks --case=snake --with-expecter
type ReviewRepositoryInterface interface {
	generic.CRUDRepository[models.Review]
	GetByReviewerId(reviewerId uint) ([]models.Review, error)
	GetByRevieweeId(reviewedId uint) ([]models.Review, error)
	GetByReviewerIdAndRevieweeId(reviewerId uint, reviewedId uint) (*models.Review, error)
	GetFiltered(filter *ReviewFilter) ([]models.Review, *pagination.PaginationResponse, error)
	GetAverageRatingByRevieweeId(revieweeId uint) (float64, error)
	GetFrequencyOfRatingByRevieweeId(revieweeId uint) (map[int]int, error)
}

type ReviewRepository struct {
	repository *generic.GormRepository[models.Review]
}

func NewReviewRepository(dbHandle *gorm.DB) ReviewRepositoryInterface {
	return &ReviewRepository{repository: generic.GetGormRepository[models.Review](dbHandle)}
}

func (repo *ReviewRepository) Create(review *models.Review) error {
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

func (repo *ReviewRepository) GetAll() ([]models.Review, error) {
	db := repo.repository.DB
	var reviews []models.Review
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

func (repo *ReviewRepository) GetById(id uint) (*models.Review, error) {
	db := repo.repository.DB
	var review models.Review
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

func (repo *ReviewRepository) GetFiltered(filter *ReviewFilter) ([]models.Review, *pagination.PaginationResponse, error) {
	query := repo.buildBaseQuery().
		Joins("JOIN users as reviewer ON reviews.reviewer_id = reviewer.id").
		Joins("JOIN users as reviewee ON reviews.reviewee_id = reviewee.id")
	query, err := filter.ApplyReviewFilters(query)
	if err != nil {
		return nil, nil, err
	}
	reviews, paginationResponse, err := pagination.PaginateResults[models.Review](&filter.Pagination, query)
	if err != nil {
		return nil, nil, err
	}
	return reviews, paginationResponse, nil
}

func (repo *ReviewRepository) Update(review *models.Review) error {
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

func (repo *ReviewRepository) GetByReviewerId(reviewerId uint) ([]models.Review, error) {
	var reviews []models.Review
	err := repo.repository.
		DB.
		Where("reviewer_id = ?", reviewerId).
		Preload("Reviewee").
		Preload("Reviewer").
		Find(&reviews).
		Error
	return reviews, err
}

func (repo *ReviewRepository) GetByRevieweeId(reviewedId uint) ([]models.Review, error) {
	var reviews []models.Review
	err := repo.repository.
		DB.
		Where("reviewee_id = ?", reviewedId).
		Preload("Reviewee").
		Preload("Reviewer").
		Find(&reviews).
		Error
	return reviews, err
}

func (repo *ReviewRepository) GetByReviewerIdAndRevieweeId(reviewerId uint, reviewedId uint) (*models.Review, error) {
	var review models.Review
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
		Model(&models.Review{}).
		Select("AVG(rating)").
		Where("reviewee_id = ?", revieweeId).
		Scan(&average).
		Error
	if err != nil {
		return 0, errors.New("no reviews found")
	}
	roundedAverage := math.Round(average*100) / 100
	return roundedAverage, nil
}

func (repo *ReviewRepository) buildBaseQuery() *gorm.DB {
	db := repo.repository.DB
	query := db.
		Preload("Reviewer").
		Preload("Reviewee")
	return query
}

func (repo *ReviewRepository) GetFrequencyOfRatingByRevieweeId(revieweeId uint) (map[int]int, error) {
	freqMap := repo.prepareFreqMap()
	raw, err := repo.getFrequencies(revieweeId)
	if err != nil {
		return freqMap, err
	}
	total, err := repo.getTotalReviews(revieweeId)
	if err != nil {
		return freqMap, err
	}

	for _, rf := range raw {
		freqMap[rf.Rating] = rf.Frequency
	}
	if total > 0 {
		for i := 1; i <= 5; i++ {
			freqMap[i] = int(float64(freqMap[i]) / float64(total) * 100)
		}
	}
	return freqMap, nil
}

func (repo *ReviewRepository) prepareFreqMap() map[int]int {
	freqMap := make(map[int]int, 5)
	for i := 1; i <= 5; i++ {
		freqMap[i] = 0
	}
	return freqMap
}

func (repo *ReviewRepository) getFrequencies(revieweeId uint) ([]RatingFrequency, error) {
	var frequencies []RatingFrequency
	err := repo.repository.
		DB.
		Model(&models.Review{}).
		Select("rating, COUNT(*) AS frequency").
		Where("reviewee_id = ?", revieweeId).
		Group("rating").
		Scan(&frequencies).
		Error
	if err != nil {
		return nil, err
	}
	return frequencies, nil
}

func (repo *ReviewRepository) getTotalReviews(revieweeId uint) (int64, error) {
	var total int64
	err := repo.repository.
		DB.
		Model(&models.Review{}).
		Where("reviewee_id = ?", revieweeId).
		Count(&total).
		Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
