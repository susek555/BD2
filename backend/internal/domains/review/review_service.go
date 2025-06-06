package review

import (
	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
)

type ReviewServiceInterface interface {
	Create(userID uint, review *CreateReviewDTO) (*RetrieveReviewDTO, error)
	GetAll() ([]RetrieveReviewDTO, error)
	GetByID(id uint) (*RetrieveReviewDTO, error)
	Update(userID uint, review *UpdateReviewDTO) (*RetrieveReviewDTO, error)
	Delete(userID, id uint) error
	GetByReviewerID(reviewerID uint) ([]RetrieveReviewDTO, error)
	GetByRevieweeID(reviewedID uint) ([]RetrieveReviewDTO, error)
	GetByReviewerIDAndRevieweeID(reviewerID uint, revieweeID uint) (*RetrieveReviewDTO, error)
	GetFiltered(filter *ReviewFilter) (*RetrieveReviewsWithPagination, error)
	GetAverageRatingByRevieweeID(revieweeID uint) (float64, error)
	GetFrequencyOfRatingByRevieweeID(revieweeID uint) (map[int]int, error)
}

type ReviewService struct {
	Repo ReviewRepositoryInterface
}

func NewReviewService(repo ReviewRepositoryInterface) ReviewServiceInterface {
	return &ReviewService{
		Repo: repo,
	}
}

func (service *ReviewService) Create(userID uint, review *CreateReviewDTO) (*RetrieveReviewDTO, error) {
	reviewObj, err := review.MapToObject(userID)
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

func (service *ReviewService) GetByID(id uint) (*RetrieveReviewDTO, error) {
	review, err := service.Repo.GetByID(id)
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

func (service *ReviewService) GetAverageRatingByRevieweeID(revieweeID uint) (float64, error) {
	averageRating, err := service.Repo.GetAverageRatingByRevieweeID(revieweeID)
	if err != nil {
		return 0, err
	}
	if averageRating == 0 {
		return 0, ErrNoReviewsFound
	}
	return averageRating, nil
}

func (service *ReviewService) Update(reviewerID uint, review *UpdateReviewDTO) (*RetrieveReviewDTO, error) {
	revieweeID, err := service.getRevieweeID(review.ID)
	if err != nil {
		return nil, err
	}
	_, err = service.Repo.GetByReviewerIDAndRevieweeID(reviewerID, revieweeID)
	if err != nil {
		return nil, ErrNotReviewer
	}
	reviewObj, err := review.MapToObject(reviewerID, revieweeID)
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

func (service *ReviewService) Delete(userID, id uint) error {
	review, err := service.Repo.GetByID(id)
	if err != nil {
		return err
	}
	if review.ReviewerID != userID {
		return ErrNotReviewer
	}
	return service.Repo.Delete(id)
}

func (service *ReviewService) GetByReviewerID(reviewerID uint) ([]RetrieveReviewDTO, error) {
	reviews, err := service.Repo.GetByReviewerID(reviewerID)
	if err != nil {
		return nil, err
	}
	reviewsDTO := mapping.MapSliceToDTOs(reviews, MapToDTO)
	return reviewsDTO, nil
}

func (service *ReviewService) GetByRevieweeID(reviewedID uint) ([]RetrieveReviewDTO, error) {
	reviews, err := service.Repo.GetByRevieweeID(reviewedID)
	if err != nil {
		return nil, err
	}
	reviewsDTO := mapping.MapSliceToDTOs(reviews, MapToDTO)
	return reviewsDTO, nil
}

func (service *ReviewService) GetByReviewerIDAndRevieweeID(reviewerID uint, reviewedID uint) (*RetrieveReviewDTO, error) {
	review, err := service.Repo.GetByReviewerIDAndRevieweeID(reviewerID, reviewedID)
	if err != nil {
		return nil, err
	}
	reviewDTO := MapToDTO(review)
	return reviewDTO, nil
}

func (service *ReviewService) getRevieweeID(reviewID uint) (uint, error) {
	review, err := service.Repo.GetByID(reviewID)
	if err != nil {
		return 0, err
	}
	return review.RevieweeID, nil
}

func (service *ReviewService) GetFrequencyOfRatingByRevieweeID(revieweeID uint) (map[int]int, error) {
	frequency, err := service.Repo.GetFrequencyOfRatingByRevieweeID(revieweeID)
	if err != nil {
		return nil, err
	}
	if len(frequency) == 0 {
		return nil, ErrNoReviewsFound
	}
	return frequency, nil
}
