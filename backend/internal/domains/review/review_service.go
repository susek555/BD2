package review

import (
	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
)

type ReviewServiceInterface interface {
	Create(userId uint, review *CreateReviewDTO) (*RetrieveReviewDTO, error)
	GetAll() ([]RetrieveReviewDTO, error)
	GetById(id uint) (*RetrieveReviewDTO, error)
	Update(userId uint, review *UpdateReviewDTO) (*RetrieveReviewDTO, error)
	Delete(userId, id uint) error
	GetByReviewerId(reviewerId uint) ([]RetrieveReviewDTO, error)
	GetByRevieweeId(reviewedId uint) ([]RetrieveReviewDTO, error)
	GetByReviewerIdAndRevieweeId(reviewerId uint, revieweeId uint) (*RetrieveReviewDTO, error)
	GetFiltered(filter *ReviewFilter) (*RetrieveReviewsWithPagination, error)
	GetAverageRatingByRevieweeId(revieweeId uint) (float64, error)
	GetFrequencyOfRatingByRevieweeId(revieweeId uint) (map[int]int, error)
}

type ReviewService struct {
	Repo ReviewRepositoryInterface
}

func NewReviewService(repo ReviewRepositoryInterface) ReviewServiceInterface {
	return &ReviewService{
		Repo: repo,
	}
}

func (service *ReviewService) Create(userId uint, review *CreateReviewDTO) (*RetrieveReviewDTO, error) {
	reviewObj, err := review.MapToObject(userId)
	if err != nil {
		return nil, err
	}
	err = service.Repo.Create(reviewObj)
	if err != nil {
		return nil, err
	}
	reviewDTO := MapToDTO(reviewObj)
	return reviewDTO, nil
}

func (service *ReviewService) GetAll() ([]RetrieveReviewDTO, error) {
	reviews, err := service.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	reviewsDTO := mapping.MapSliceToDTOs(reviews, MapToDTO)
	return reviewsDTO, nil
}

func (service *ReviewService) GetById(id uint) (*RetrieveReviewDTO, error) {
	review, err := service.Repo.GetById(id)
	if err != nil {
		return nil, ErrNoReviewFound
	}
	reviewDTO := MapToDTO(review)
	return reviewDTO, nil
}

func (service *ReviewService) GetFiltered(filter *ReviewFilter) (*RetrieveReviewsWithPagination, error) {
	reviews, pagResponse, err := service.Repo.GetFiltered(filter)
	if err != nil {
		return nil, err
	}
	reviewsDTO := mapping.MapSliceToDTOs(reviews, MapToDTO)
	return &RetrieveReviewsWithPagination{
		Reviews:            reviewsDTO,
		PaginationResponse: pagResponse,
	}, nil
}

func (service *ReviewService) GetAverageRatingByRevieweeId(revieweeId uint) (float64, error) {
	averageRating, err := service.Repo.GetAverageRatingByRevieweeId(revieweeId)
	if err != nil {
		return 0, err
	}
	if averageRating == 0 {
		return 0, ErrNoReviewsFound
	}
	return averageRating, nil
}

func (service *ReviewService) Update(reviewerId uint, review *UpdateReviewDTO) (*RetrieveReviewDTO, error) {
	revieweeId, err := service.getRevieweeId(review.ID)
	if err != nil {
		return nil, err
	}
	_, err = service.Repo.GetByReviewerIdAndRevieweeId(reviewerId, revieweeId)
	if err != nil {
		return nil, ErrNotReviewer
	}
	reviewObj, err := review.MapToObject(reviewerId, revieweeId)
	if err != nil {
		return nil, err
	}
	err = service.Repo.Update(reviewObj)
	if err != nil {
		return nil, err
	}
	reviewDTO := MapToDTO(reviewObj)
	return reviewDTO, nil
}

func (service *ReviewService) Delete(userId, id uint) error {
	review, err := service.Repo.GetById(id)
	if err != nil {
		return err
	}
	if review.ReviewerID != userId {
		return ErrNotReviewer
	}
	return service.Repo.Delete(id)
}

func (service *ReviewService) GetByReviewerId(reviewerId uint) ([]RetrieveReviewDTO, error) {
	reviews, err := service.Repo.GetByReviewerId(reviewerId)
	if err != nil {
		return nil, err
	}
	reviewsDTO := mapping.MapSliceToDTOs(reviews, MapToDTO)
	return reviewsDTO, nil
}

func (service *ReviewService) GetByRevieweeId(reviewedId uint) ([]RetrieveReviewDTO, error) {
	reviews, err := service.Repo.GetByRevieweeId(reviewedId)
	if err != nil {
		return nil, err
	}
	reviewsDTO := mapping.MapSliceToDTOs(reviews, MapToDTO)
	return reviewsDTO, nil
}

func (service *ReviewService) GetByReviewerIdAndRevieweeId(reviewerId uint, reviewedId uint) (*RetrieveReviewDTO, error) {
	review, err := service.Repo.GetByReviewerIdAndRevieweeId(reviewerId, reviewedId)
	if err != nil {
		return nil, err
	}
	reviewDTO := MapToDTO(review)
	return reviewDTO, nil
}

func (service *ReviewService) getRevieweeId(reviewId uint) (uint, error) {
	review, err := service.Repo.GetById(reviewId)
	if err != nil {
		return 0, err
	}
	return review.RevieweeId, nil
}

func (service *ReviewService) GetFrequencyOfRatingByRevieweeId(revieweeId uint) (map[int]int, error) {
	frequency, err := service.Repo.GetFrequencyOfRatingByRevieweeId(revieweeId)
	if err != nil {
		return nil, err
	}
	if len(frequency) == 0 {
		return nil, ErrNoReviewsFound
	}
	return frequency, nil
}
